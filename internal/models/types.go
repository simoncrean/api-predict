package models

import "time"

// SystemSpec represents a user's system specifications
type SystemSpec struct {
	CPUCores    int    `json:"cpu_cores" binding:"required,min=1,max=64"`
	RAMGB       int    `json:"ram_gb" binding:"required,min=1,max=128"`
	StorageGB   int    `json:"storage_gb" binding:"required,min=32,max=8192"`
	HasSSD      bool   `json:"has_ssd"`
	HasGPU      bool   `json:"has_gpu"`
	GPUVRAMGB   int    `json:"gpu_vram_gb" binding:"min=0,max=48"`
	NetworkMbps int    `json:"network_mbps" binding:"required,min=1,max=10000"`
	OS          string `json:"os" binding:"required,oneof=Windows Linux macOS"`
}

// DePINProject represents a DePIN project specification
type DePINProject struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	NodeType         string `json:"node_type"`
	CPUCoresMin      int    `json:"cpu_cores_min"`
	RAMGBMin         int    `json:"ram_gb_min"`
	RAMGBRecommended int    `json:"ram_gb_recommended"`
	StorageGBMin     int    `json:"storage_gb_min"`
	StorageType      string `json:"storage_type"` // "SSD", "Any"
	GPURequired      bool   `json:"gpu_required"`
	GPUVRAMGBMin     int    `json:"gpu_vram_gb_min"`
	NetworkMbpsMin   int    `json:"network_mbps_min"`
	SupportedOS      string `json:"supported_os"` // "Linux,Windows,macOS"
	EstimatedCostMin int    `json:"estimated_cost_min"`
	EstimatedCostMax int    `json:"estimated_cost_max"`
	CostCategory     string `json:"cost_category"`
	HomeFriendly     bool   `json:"home_friendly"`
	Description      string `json:"description"`
}

// CompatibilityResult represents the compatibility analysis for a single project
type CompatibilityResult struct {
	Name                string   `json:"name"`
	Compatible          bool     `json:"compatible"`
	CompatibilityScore  float64  `json:"compatibility_score"`
	PerformanceRating   string   `json:"performance_rating"`
	EstimatedCost       string   `json:"estimated_cost"`
	MissingRequirements []string `json:"missing_requirements"`
	RecommendedUpgrades []string `json:"recommended_upgrades"`
	Warnings            []string `json:"warnings,omitempty"`
}

// PredictionRequest represents the API request for compatibility prediction
type PredictionRequest struct {
	System SystemSpec `json:"system" binding:"required"`
}

// PredictionResponse represents the API response with compatibility results
type PredictionResponse struct {
	CompatibleProjects   []CompatibilityResult `json:"compatible_projects"`
	IncompatibleProjects []CompatibilityResult `json:"incompatible_projects"`
	Summary              PredictionSummary     `json:"summary"`
	Recommendations      []string              `json:"recommendations"`
	GeneratedAt          time.Time             `json:"generated_at"`
}

// PredictionSummary provides overview statistics
type PredictionSummary struct {
	TotalProjects     int     `json:"total_projects"`
	CompatibleCount   int     `json:"compatible_count"`
	IncompatibleCount int     `json:"incompatible_count"`
	CompatibilityRate float64 `json:"compatibility_rate"`
	AverageScore      float64 `json:"average_score"`
	SystemRating      string  `json:"system_rating"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status         string    `json:"status"`
	Version        string    `json:"version"`
	ProjectsLoaded int       `json:"projects_loaded"`
	Uptime         string    `json:"uptime"`
	Timestamp      time.Time `json:"timestamp"`
}

// ProjectsResponse represents the response for listing all projects
type ProjectsResponse struct {
	Projects []DePINProject `json:"projects"`
	Total    int            `json:"total"`
	Summary  ProjectSummary `json:"summary"`
}

// ProjectSummary provides statistics about loaded projects
type ProjectSummary struct {
	ByType         map[string]int `json:"by_type"`
	ByCostCategory map[string]int `json:"by_cost_category"`
	HomeFriendly   int            `json:"home_friendly"`
	GPURequired    int            `json:"gpu_required"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string    `json:"error"`
	Message string    `json:"message,omitempty"`
	Code    int       `json:"code"`
	Time    time.Time `json:"timestamp"`
}

// Performance ratings
const (
	RatingExcellent = "Excellent"
	RatingGood      = "Good"
	RatingFair      = "Fair"
	RatingPoor      = "Poor"
)

// System categories for rating
const (
	SystemEntry    = "Entry Level"
	SystemMidRange = "Mid-Range"
	SystemHighEnd  = "High-End"
	SystemExtreme  = "Extreme"
)

// Compatibility thresholds
const (
	ScoreExcellent = 0.9
	ScoreGood      = 0.7
	ScoreFair      = 0.5
	ScorePoor      = 0.0
)

// GetPerformanceRating returns performance rating based on score
func GetPerformanceRating(score float64) string {
	switch {
	case score >= ScoreExcellent:
		return RatingExcellent
	case score >= ScoreGood:
		return RatingGood
	case score >= ScoreFair:
		return RatingFair
	default:
		return RatingPoor
	}
}

// GetSystemRating categorizes system based on specifications
func GetSystemRating(spec SystemSpec) string {
	score := 0

	// CPU scoring
	if spec.CPUCores >= 12 {
		score += 3
	} else if spec.CPUCores >= 8 {
		score += 2
	} else if spec.CPUCores >= 4 {
		score += 1
	}

	// RAM scoring
	if spec.RAMGB >= 32 {
		score += 3
	} else if spec.RAMGB >= 16 {
		score += 2
	} else if spec.RAMGB >= 8 {
		score += 1
	}

	// GPU scoring
	if spec.HasGPU && spec.GPUVRAMGB >= 12 {
		score += 3
	} else if spec.HasGPU && spec.GPUVRAMGB >= 6 {
		score += 2
	} else if spec.HasGPU {
		score += 1
	}

	// Storage scoring
	if spec.HasSSD && spec.StorageGB >= 1000 {
		score += 2
	} else if spec.HasSSD || spec.StorageGB >= 500 {
		score += 1
	}

	// Network scoring
	if spec.NetworkMbps >= 500 {
		score += 2
	} else if spec.NetworkMbps >= 100 {
		score += 1
	}

	// Categorize based on total score
	switch {
	case score >= 12:
		return SystemExtreme
	case score >= 8:
		return SystemHighEnd
	case score >= 5:
		return SystemMidRange
	default:
		return SystemEntry
	}
}
