package api

import (
	"net/http"
	"time"

	"github/simoncrean/api-predict/internal/models"
	"github/simoncrean/api-predict/internal/service"

	"github.com/gin-gonic/gin"
)

// Handlers contains all HTTP request handlers
type Handlers struct {
	compatibilityService *service.CompatibilityService
}

// NewHandlers creates a new handlers instance
func NewHandlers(compatibilityService *service.CompatibilityService) *Handlers {
	return &Handlers{
		compatibilityService: compatibilityService,
	}
}

// PredictCompatibility handles DePIN compatibility prediction requests
func (h *Handlers) PredictCompatibility(c *gin.Context) {
	var request models.PredictionRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
			Time:    time.Now(),
		})
		return
	}

	// Validate system specifications
	if err := validateSystemSpec(request.System); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid system specifications",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
			Time:    time.Now(),
		})
		return
	}

	// Perform compatibility prediction
	result, err := h.compatibilityService.PredictCompatibility(request.System)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Prediction failed",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
			Time:    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(c *gin.Context) {
	projects := h.compatibilityService.GetProjects()
	uptime := h.compatibilityService.GetUptime()

	health := models.HealthResponse{
		Status:         "healthy",
		Version:        "1.0.0",
		ProjectsLoaded: len(projects),
		Uptime:         uptime.String(),
		Timestamp:      time.Now(),
	}

	c.JSON(http.StatusOK, health)
}

// ListProjects handles requests for listing all DePIN projects
func (h *Handlers) ListProjects(c *gin.Context) {
	projects := h.compatibilityService.GetProjects()
	summary := h.compatibilityService.GetProjectSummary()

	response := models.ProjectsResponse{
		Projects: projects,
		Total:    len(projects),
		Summary:  summary,
	}

	c.JSON(http.StatusOK, response)
}

// APIDocs serves API documentation
func (h *Handlers) APIDocs(c *gin.Context) {
	docs := gin.H{
		"service":     "DePIN Compatibility API",
		"version":     "1.0.0",
		"description": "Predicts DePIN compatibility based on consumer system specifications",
		"endpoints": gin.H{
			"POST /api/v1/predict": gin.H{
				"description": "Predict DePIN compatibility for a system",
				"example_request": gin.H{
					"system": gin.H{
						"cpu_cores":    8,
						"ram_gb":       16,
						"storage_gb":   512,
						"has_ssd":      true,
						"has_gpu":      true,
						"gpu_vram_gb":  8,
						"network_mbps": 100,
						"os":           "Windows",
					},
				},
			},
			"GET /api/v1/health": gin.H{
				"description": "Service health check",
			},
			"GET /api/v1/projects": gin.H{
				"description": "List all DePIN projects",
			},
			"GET /api/v1/metrics": gin.H{
				"description": "Service metrics",
			},
		},
		"system_requirements": gin.H{
			"cpu_cores":    "Number of CPU cores (1-64)",
			"ram_gb":       "RAM in GB (1-128)",
			"storage_gb":   "Storage in GB (32-8192)",
			"has_ssd":      "Boolean - SSD storage",
			"has_gpu":      "Boolean - Dedicated GPU",
			"gpu_vram_gb":  "GPU VRAM in GB (0-48)",
			"network_mbps": "Network speed in Mbps (1-10000)",
			"os":           "Operating system (Windows/Linux/macOS)",
		},
		"compatibility_scores": gin.H{
			"excellent": "0.9 - 1.0",
			"good":      "0.7 - 0.89",
			"fair":      "0.5 - 0.69",
			"poor":      "0.0 - 0.49",
		},
	}

	c.JSON(http.StatusOK, docs)
}

// Metrics handles metrics requests (simplified Prometheus-style metrics)
func (h *Handlers) Metrics(c *gin.Context) {
	projects := h.compatibilityService.GetProjects()
	summary := h.compatibilityService.GetProjectSummary()
	uptime := h.compatibilityService.GetUptime()

	metrics := gin.H{
		"service_info": gin.H{
			"name":           "depin_compatibility_api",
			"version":        "1.0.0",
			"uptime_seconds": uptime.Seconds(),
		},
		"projects_loaded_total":  len(projects),
		"projects_by_type":       summary.ByType,
		"projects_by_cost":       summary.ByCostCategory,
		"projects_home_friendly": summary.HomeFriendly,
		"projects_gpu_required":  summary.GPURequired,
		"timestamp":              time.Now().Unix(),
	}

	c.JSON(http.StatusOK, metrics)
}

// validateSystemSpec performs additional validation on system specifications
func validateSystemSpec(spec models.SystemSpec) error {
	// Custom validation logic can be added here
	// For example, logical consistency checks

	// If GPU is claimed but no VRAM specified
	if spec.HasGPU && spec.GPUVRAMGB == 0 {
		// This might be valid for very old GPUs, so just a warning
	}

	// If no GPU but VRAM specified
	if !spec.HasGPU && spec.GPUVRAMGB > 0 {
		// This might indicate integrated graphics
	}

	return nil
}
