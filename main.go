package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-predict/internal/api"
	"api-predict/internal/data"
	"api-predict/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort     = "8080"
	defaultHost     = "0.0.0.0"
	defaultDataPath = "./data/depin_specs.csv" // Will use depin_specifications_final.csv if available
)

func main() {
	// Load configuration from environment
	config := loadConfig()

	// Initialize data loader
	dataLoader := data.NewLoader(config.DataPath)
	depinProjects, err := dataLoader.LoadDePINSpecs()
	if err != nil {
		log.Fatalf("Failed to load DePIN specifications: %v", err)
	}

	log.Printf("Loaded %d DePIN projects", len(depinProjects))

	// Initialize services
	compatibilityService := service.NewCompatibilityService(depinProjects)

	// Initialize API handlers
	handlers := api.NewHandlers(compatibilityService)

	// Setup router
	router := setupRouter(handlers)

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 DePIN Compatibility API starting on %s:%s", config.Host, config.Port)
		log.Printf("📊 Health check: http://localhost:%s/api/v1/health", config.Port)
		log.Printf("📖 API docs: http://localhost:%s/api/v1/docs", config.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	} else {
		log.Println("✅ Server shutdown complete")
	}
}

// Config holds application configuration
type Config struct {
	Port     string
	Host     string
	DataPath string
	LogLevel string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	return &Config{
		Port:     getEnv("PORT", defaultPort),
		Host:     getEnv("HOST", defaultHost),
		DataPath: getEnv("DATA_PATH", defaultDataPath),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// setupRouter configures the HTTP router
func setupRouter(handlers *api.Handlers) *gin.Engine {
	// Set gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(api.CORSMiddleware())
	router.Use(api.RateLimitMiddleware())

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Core endpoints
		v1.POST("/predict", handlers.PredictCompatibility)
		v1.GET("/health", handlers.HealthCheck)
		v1.GET("/projects", handlers.ListProjects)

		// Utility endpoints
		v1.GET("/docs", handlers.APIDocs)
		v1.GET("/metrics", handlers.Metrics)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "DePIN Resource Predict API",
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
			"health":  "/api/v1/health",
		})
	})

	return router
}
