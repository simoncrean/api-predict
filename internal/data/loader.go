package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/simoncrean/api-predict/internal/models"
)

// Loader handles loading DePIN project data from CSV files
type Loader struct {
	filePath string
}

// NewLoader creates a new data loader
func NewLoader(filePath string) *Loader {
	return &Loader{
		filePath: filePath,
	}
}

// LoadDePINSpecs loads DePIN project specifications from CSV file
func (l *Loader) LoadDePINSpecs() ([]models.DePINProject, error) {
	file, err := os.Open(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %w", l.filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header to create field mapping
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	fieldMap := createFieldMap(header)

	var projects []models.DePINProject

	// Read data rows
	lineNumber := 2 // Start from line 2 (after header)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV line %d: %w", lineNumber, err)
		}

		project, err := l.parseProjectRecord(record, fieldMap, lineNumber)
		if err != nil {
			// Log warning but continue processing
			fmt.Printf("Warning: Failed to parse line %d: %v\n", lineNumber, err)
			lineNumber++
			continue
		}

		projects = append(projects, project)
		lineNumber++
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("no valid projects found in CSV file")
	}

	return projects, nil
}

// createFieldMap creates a mapping from field names to column indices
func createFieldMap(header []string) map[string]int {
	fieldMap := make(map[string]int)
	for i, field := range header {
		// Normalize field names (remove spaces, convert to lowercase)
		normalizedField := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(field), " ", "_"))
		fieldMap[normalizedField] = i
	}
	return fieldMap
}

// parseProjectRecord parses a single CSV record into a DePINProject
func (l *Loader) parseProjectRecord(record []string, fieldMap map[string]int, lineNumber int) (models.DePINProject, error) {
	project := models.DePINProject{}

	// Project name (required)
	project.Name = getStringField(record, fieldMap, "project_name", "name")
	if project.Name == "" {
		return project, fmt.Errorf("project name is required")
	}

	// Project type
	project.Type = getStringField(record, fieldMap, "project_type", "type")

	// Node type
	project.NodeType = getStringField(record, fieldMap, "node_type")

	// CPU requirements
	project.CPUCoresMin = getIntField(record, fieldMap, "cpu_cores_min")

	// RAM requirements
	project.RAMGBMin = getIntField(record, fieldMap, "ram_gb_min", "ram_min_gb")
	project.RAMGBRecommended = getIntField(record, fieldMap, "ram_gb_recommended", "ram_recommended_gb")

	// Storage requirements
	project.StorageGBMin = getIntField(record, fieldMap, "storage_gb_min", "storage_min_gb")
	project.StorageType = getStringField(record, fieldMap, "storage_type")

	// GPU requirements
	project.GPURequired = getBoolField(record, fieldMap, "gpu_required")
	project.GPUVRAMGBMin = getIntField(record, fieldMap, "gpu_vram_gb_min", "gpu_vram_min_gb")

	// Network requirements
	project.NetworkMbpsMin = getIntField(record, fieldMap, "network_speed_mbps_min", "network_mbps_min")

	// Supported OS
	project.SupportedOS = getStringField(record, fieldMap, "supported_os", "os_support")

	// Cost estimates
	project.EstimatedCostMin = getIntField(record, fieldMap, "estimated_monthly_cost_usd_min", "cost_min")
	project.EstimatedCostMax = getIntField(record, fieldMap, "estimated_monthly_cost_usd_max", "cost_max")
	project.CostCategory = getStringField(record, fieldMap, "cost_category")

	// Home friendly
	project.HomeFriendly = getBoolField(record, fieldMap, "home_friendly")

	// Description
	project.Description = getStringField(record, fieldMap, "description", "additional_requirements")

	// Validate required fields
	if err := l.validateProject(project); err != nil {
		return project, fmt.Errorf("validation failed: %w", err)
	}

	return project, nil
}

// validateProject validates that a project has required fields and sensible values
func (l *Loader) validateProject(project models.DePINProject) error {
	if project.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if project.CPUCoresMin < 0 || project.CPUCoresMin > 64 {
		return fmt.Errorf("invalid CPU cores minimum: %d", project.CPUCoresMin)
	}

	if project.RAMGBMin < 0 || project.RAMGBMin > 1024 {
		return fmt.Errorf("invalid RAM minimum: %d", project.RAMGBMin)
	}

	if project.StorageGBMin < 0 || project.StorageGBMin > 100000 {
		return fmt.Errorf("invalid storage minimum: %d", project.StorageGBMin)
	}

	if project.NetworkMbpsMin < 0 || project.NetworkMbpsMin > 100000 {
		return fmt.Errorf("invalid network speed minimum: %d", project.NetworkMbpsMin)
	}

	// Set defaults for missing optional fields
	if project.Type == "" {
		project.Type = "Unknown"
	}

	if project.NodeType == "" {
		project.NodeType = "Standard"
	}

	if project.StorageType == "" {
		project.StorageType = "Any"
	}

	if project.CostCategory == "" {
		if project.EstimatedCostMax <= 20 {
			project.CostCategory = "Low"
		} else if project.EstimatedCostMax <= 100 {
			project.CostCategory = "Medium"
		} else {
			project.CostCategory = "High"
		}
	}

	if project.SupportedOS == "" {
		project.SupportedOS = "Linux,Windows,macOS"
	}

	return nil
}

// Helper functions for extracting fields from CSV records

func getStringField(record []string, fieldMap map[string]int, fieldNames ...string) string {
	for _, fieldName := range fieldNames {
		if idx, ok := fieldMap[fieldName]; ok && idx < len(record) {
			value := strings.TrimSpace(record[idx])
			if value != "" {
				return value
			}
		}
	}
	return ""
}

func getIntField(record []string, fieldMap map[string]int, fieldNames ...string) int {
	for _, fieldName := range fieldNames {
		if idx, ok := fieldMap[fieldName]; ok && idx < len(record) {
			value := strings.TrimSpace(record[idx])
			if value != "" {
				if intVal, err := strconv.Atoi(value); err == nil {
					return intVal
				}
			}
		}
	}
	return 0
}

func getBoolField(record []string, fieldMap map[string]int, fieldNames ...string) bool {
	for _, fieldName := range fieldNames {
		if idx, ok := fieldMap[fieldName]; ok && idx < len(record) {
			value := strings.ToUpper(strings.TrimSpace(record[idx]))
			return value == "TRUE" || value == "1" || value == "YES" || value == "Y"
		}
	}
	return false
}
