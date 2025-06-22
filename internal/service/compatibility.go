package service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/simoncrean/api-predict/internal/models"
)

// CompatibilityService handles DePIN compatibility analysis
type CompatibilityService struct {
	projects  []models.DePINProject
	startTime time.Time
}

// NewCompatibilityService creates a new compatibility service
func NewCompatibilityService(projects []models.DePINProject) *CompatibilityService {
	return &CompatibilityService{
		projects:  projects,
		startTime: time.Now(),
	}
}

// PredictCompatibility analyzes system compatibility with all DePIN projects
func (s *CompatibilityService) PredictCompatibility(system models.SystemSpec) (*models.PredictionResponse, error) {
	var compatible []models.CompatibilityResult
	var incompatible []models.CompatibilityResult
	totalScore := 0.0

	for _, project := range s.projects {
		result := s.analyzeProjectCompatibility(system, project)

		if result.Compatible {
			compatible = append(compatible, result)
		} else {
			incompatible = append(incompatible, result)
		}

		totalScore += result.CompatibilityScore
	}

	// Sort results by compatibility score (descending)
	sort.Slice(compatible, func(i, j int) bool {
		return compatible[i].CompatibilityScore > compatible[j].CompatibilityScore
	})

	sort.Slice(incompatible, func(i, j int) bool {
		return incompatible[i].CompatibilityScore > incompatible[j].CompatibilityScore
	})

	// Calculate summary statistics
	summary := models.PredictionSummary{
		TotalProjects:     len(s.projects),
		CompatibleCount:   len(compatible),
		IncompatibleCount: len(incompatible),
		CompatibilityRate: float64(len(compatible)) / float64(len(s.projects)) * 100,
		AverageScore:      totalScore / float64(len(s.projects)),
		SystemRating:      models.GetSystemRating(system),
	}

	// Generate recommendations
	recommendations := s.generateRecommendations(system, compatible, incompatible)

	return &models.PredictionResponse{
		CompatibleProjects:   compatible,
		IncompatibleProjects: incompatible,
		Summary:              summary,
		Recommendations:      recommendations,
		GeneratedAt:          time.Now(),
	}, nil
}

// analyzeProjectCompatibility performs detailed compatibility analysis for a single project
func (s *CompatibilityService) analyzeProjectCompatibility(system models.SystemSpec, project models.DePINProject) models.CompatibilityResult {
	result := models.CompatibilityResult{
		Name:                project.Name,
		Compatible:          true,
		CompatibilityScore:  1.0,
		PerformanceRating:   models.RatingExcellent,
		EstimatedCost:       fmt.Sprintf("$%d-$%d/month", project.EstimatedCostMin, project.EstimatedCostMax),
		MissingRequirements: []string{},
		RecommendedUpgrades: []string{},
		Warnings:            []string{},
	}

	score := 1.0

	// Check CPU requirements
	if system.CPUCores < project.CPUCoresMin {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements,
			fmt.Sprintf("CPU cores: need %d, have %d", project.CPUCoresMin, system.CPUCores))
		score -= 0.3
	}

	// Check RAM requirements
	if system.RAMGB < project.RAMGBMin {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements,
			fmt.Sprintf("RAM: need %dGB, have %dGB", project.RAMGBMin, system.RAMGB))
		score -= 0.3
	} else if system.RAMGB < project.RAMGBRecommended {
		result.RecommendedUpgrades = append(result.RecommendedUpgrades,
			fmt.Sprintf("RAM upgrade to %dGB recommended for optimal performance", project.RAMGBRecommended))
		score -= 0.1
	}

	// Check storage requirements
	if system.StorageGB < project.StorageGBMin {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements,
			fmt.Sprintf("Storage: need %dGB, have %dGB", project.StorageGBMin, system.StorageGB))
		score -= 0.2
	}

	// Check SSD requirement
	if project.StorageType == "SSD" && !system.HasSSD {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements, "SSD storage required")
		score -= 0.25
	} else if project.StorageType == "SSD" && system.HasSSD {
		// Bonus for having SSD when recommended
		score += 0.05
	}

	// Check GPU requirements
	if project.GPURequired && !system.HasGPU {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements, "Dedicated GPU required")
		score -= 0.4
	} else if project.GPUVRAMGBMin > 0 && system.GPUVRAMGB < project.GPUVRAMGBMin {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements,
			fmt.Sprintf("GPU VRAM: need %dGB, have %dGB", project.GPUVRAMGBMin, system.GPUVRAMGB))
		score -= 0.3
	}

	// Check network speed
	if system.NetworkMbps < project.NetworkMbpsMin {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements,
			fmt.Sprintf("Network speed: need %dMbps, have %dMbps", project.NetworkMbpsMin, system.NetworkMbps))
		score -= 0.2
	}

	// Check OS compatibility
	if !s.isOSCompatible(system.OS, project.SupportedOS) {
		result.Compatible = false
		result.MissingRequirements = append(result.MissingRequirements,
			fmt.Sprintf("OS not supported: need one of [%s], have %s", project.SupportedOS, system.OS))
		score -= 0.3
	}

	// Performance bonuses for exceeding requirements
	score += s.calculatePerformanceBonus(system, project)

	// Home-friendly check
	if !project.HomeFriendly {
		result.Warnings = append(result.Warnings, "This project may not be suitable for home use")
	}

	// Ensure score is within bounds
	score = math.Max(0.0, math.Min(1.0, score))

	result.CompatibilityScore = score
	result.PerformanceRating = models.GetPerformanceRating(score)

	return result
}

// isOSCompatible checks if the system OS is supported by the project
func (s *CompatibilityService) isOSCompatible(systemOS, supportedOS string) bool {
	if supportedOS == "" {
		return true // No restriction
	}

	supportedList := strings.Split(supportedOS, ",")
	for _, os := range supportedList {
		if strings.TrimSpace(os) == systemOS {
			return true
		}
	}
	return false
}

// calculatePerformanceBonus adds bonus points for systems that exceed requirements
func (s *CompatibilityService) calculatePerformanceBonus(system models.SystemSpec, project models.DePINProject) float64 {
	bonus := 0.0

	// CPU bonus
	if system.CPUCores > project.CPUCoresMin*2 {
		bonus += 0.05
	} else if system.CPUCores > project.CPUCoresMin {
		bonus += 0.02
	}

	// RAM bonus
	if system.RAMGB > project.RAMGBRecommended*2 {
		bonus += 0.05
	} else if system.RAMGB > project.RAMGBRecommended {
		bonus += 0.02
	}

	// Network bonus
	if system.NetworkMbps > project.NetworkMbpsMin*2 {
		bonus += 0.03
	}

	// High-end GPU bonus
	if system.HasGPU && system.GPUVRAMGB > 8 {
		bonus += 0.02
	}

	return math.Min(bonus, 0.15) // Cap bonus at 15%
}

// generateRecommendations creates personalized recommendations
func (s *CompatibilityService) generateRecommendations(system models.SystemSpec, compatible, incompatible []models.CompatibilityResult) []string {
	var recommendations []string

	compatibilityRate := float64(len(compatible)) / float64(len(s.projects))

	// Overall system assessment
	switch {
	case compatibilityRate >= 0.8:
		recommendations = append(recommendations, "🎉 Excellent! Your system is compatible with most DePIN projects.")
	case compatibilityRate >= 0.6:
		recommendations = append(recommendations, "👍 Good compatibility! Your system works well with many DePIN projects.")
	case compatibilityRate >= 0.4:
		recommendations = append(recommendations, "⚠️ Fair compatibility. Consider upgrading for better project support.")
	default:
		recommendations = append(recommendations, "📈 Limited compatibility. Upgrades recommended for better DePIN support.")
	}

	// Specific upgrade recommendations
	upgradeRecommendations := s.analyzeUpgradeNeeds(system, incompatible)
	recommendations = append(recommendations, upgradeRecommendations...)

	// Project-specific recommendations
	if len(compatible) > 0 {
		bestProjects := s.getBestProjects(compatible, 3)
		projectNames := make([]string, len(bestProjects))
		for i, project := range bestProjects {
			projectNames[i] = project.Name
		}
		recommendations = append(recommendations,
			fmt.Sprintf("🚀 Recommended projects for your system: %s", strings.Join(projectNames, ", ")))
	}

	return recommendations
}

// analyzeUpgradeNeeds suggests specific hardware upgrades
func (s *CompatibilityService) analyzeUpgradeNeeds(system models.SystemSpec, incompatible []models.CompatibilityResult) []string {
	var recommendations []string

	// Count common missing requirements
	ramIssues := 0
	cpuIssues := 0
	gpuIssues := 0
	storageIssues := 0
	networkIssues := 0

	for _, result := range incompatible {
		for _, req := range result.MissingRequirements {
			switch {
			case strings.Contains(req, "RAM"):
				ramIssues++
			case strings.Contains(req, "CPU"):
				cpuIssues++
			case strings.Contains(req, "GPU"):
				gpuIssues++
			case strings.Contains(req, "Storage") || strings.Contains(req, "SSD"):
				storageIssues++
			case strings.Contains(req, "Network"):
				networkIssues++
			}
		}
	}

	threshold := len(incompatible) / 3 // If 1/3 of projects need upgrade

	// Generate specific recommendations
	if ramIssues > threshold {
		recommendations = append(recommendations, "💾 Consider upgrading RAM for better project compatibility")
	}

	if cpuIssues > threshold {
		recommendations = append(recommendations, "🖥️ A CPU upgrade would significantly improve project support")
	}

	if gpuIssues > threshold {
		recommendations = append(recommendations, "🎮 Adding a dedicated GPU would unlock AI and compute-intensive projects")
	}

	if storageIssues > threshold {
		recommendations = append(recommendations, "💿 Consider upgrading to SSD storage or increasing capacity")
	}

	if networkIssues > threshold {
		recommendations = append(recommendations, "🌐 Faster internet connection would improve project compatibility")
	}

	return recommendations
}

// getBestProjects returns the top N compatible projects
func (s *CompatibilityService) getBestProjects(compatible []models.CompatibilityResult, n int) []models.CompatibilityResult {
	if len(compatible) <= n {
		return compatible
	}
	return compatible[:n]
}

// GetProjects returns all loaded DePIN projects
func (s *CompatibilityService) GetProjects() []models.DePINProject {
	return s.projects
}

// GetProjectSummary returns summary statistics about loaded projects
func (s *CompatibilityService) GetProjectSummary() models.ProjectSummary {
	summary := models.ProjectSummary{
		ByType:         make(map[string]int),
		ByCostCategory: make(map[string]int),
		HomeFriendly:   0,
		GPURequired:    0,
	}

	for _, project := range s.projects {
		// Count by type
		summary.ByType[project.Type]++

		// Count by cost category
		summary.ByCostCategory[project.CostCategory]++

		// Count home-friendly projects
		if project.HomeFriendly {
			summary.HomeFriendly++
		}

		// Count GPU-required projects
		if project.GPURequired {
			summary.GPURequired++
		}
	}

	return summary
}

// GetUptime returns service uptime
func (s *CompatibilityService) GetUptime() time.Duration {
	return time.Since(s.startTime)
}
