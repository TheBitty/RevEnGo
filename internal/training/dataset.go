package training

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// DatasetItem represents a single training example
type DatasetItem struct {
	Input    string                 `json:"input"`
	Output   string                 `json:"output"`
	Type     string                 `json:"type,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Dataset represents a collection of training examples
type Dataset struct {
	Items     []DatasetItem
	Path      string
	Format    string
	ItemCount int
}

// LoadDataset loads a training dataset from a directory or file
func LoadDataset(path string) (*Dataset, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access dataset: %w", err)
	}

	dataset := &Dataset{
		Path: path,
	}

	if fileInfo.IsDir() {
		// Directory format: load all .json files
		err = loadDatasetFromDir(dataset, path)
	} else {
		// File format: try to determine the format
		ext := strings.ToLower(filepath.Ext(path))

		switch ext {
		case ".json", ".jsonl":
			dataset.Format = ext[1:] // Remove the dot
			err = loadDatasetFromFile(dataset, path)
		default:
			return nil, fmt.Errorf("unsupported dataset format: %s", ext)
		}
	}

	if err != nil {
		return nil, err
	}

	dataset.ItemCount = len(dataset.Items)
	return dataset, nil
}

// loadDatasetFromDir loads a dataset from a directory containing JSON files
func loadDatasetFromDir(dataset *Dataset, dirPath string) error {
	dataset.Format = "directory"
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext != ".json" && ext != ".jsonl" {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		var fileDataset Dataset
		err := loadDatasetFromFile(&fileDataset, filePath)
		if err != nil {
			return fmt.Errorf("failed to load file %s: %w", file.Name(), err)
		}

		dataset.Items = append(dataset.Items, fileDataset.Items...)
	}

	if len(dataset.Items) == 0 {
		return fmt.Errorf("no valid dataset items found in directory")
	}

	return nil
}

// loadDatasetFromFile loads a dataset from a JSON or JSONL file
func loadDatasetFromFile(dataset *Dataset, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".jsonl" {
		// JSONL format: one JSON object per line
		lines := strings.Split(string(data), "\n")
		for lineNum, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			var item DatasetItem
			if err := json.Unmarshal([]byte(line), &item); err != nil {
				return fmt.Errorf("invalid JSON on line %d: %w", lineNum+1, err)
			}

			dataset.Items = append(dataset.Items, item)
		}
	} else {
		// Regular JSON format: can be an array or a single object
		if strings.TrimSpace(string(data))[0] == '[' {
			// JSON array format
			var items []DatasetItem
			if err := json.Unmarshal(data, &items); err != nil {
				return fmt.Errorf("invalid JSON array: %w", err)
			}
			dataset.Items = append(dataset.Items, items...)
		} else {
			// Single JSON object format
			var item DatasetItem
			if err := json.Unmarshal(data, &item); err != nil {
				return fmt.Errorf("invalid JSON object: %w", err)
			}
			dataset.Items = append(dataset.Items, item)
		}
	}

	return nil
}

// SaveDataset saves a dataset to a file
func SaveDataset(dataset *Dataset, outputPath string, format string) error {
	if format == "" {
		format = "json" // Default format
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	var data []byte
	var err error

	switch format {
	case "json":
		// Save as JSON array
		data, err = json.MarshalIndent(dataset.Items, "", "  ")
	case "jsonl":
		// Save as JSONL (one JSON object per line)
		var lines []string
		for _, item := range dataset.Items {
			itemData, err := json.Marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal item: %w", err)
			}
			lines = append(lines, string(itemData))
		}
		data = []byte(strings.Join(lines, "\n"))
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal dataset: %w", err)
	}

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Split splits a dataset into training and validation sets
func (d *Dataset) Split(trainRatio float64) (*Dataset, *Dataset) {
	trainCount := int(float64(len(d.Items)) * trainRatio)
	if trainCount <= 0 {
		trainCount = 1
	}
	if trainCount >= len(d.Items) {
		trainCount = len(d.Items) - 1
	}

	trainDataset := &Dataset{
		Path:   d.Path + "_train",
		Format: d.Format,
		Items:  d.Items[:trainCount],
	}

	valDataset := &Dataset{
		Path:   d.Path + "_val",
		Format: d.Format,
		Items:  d.Items[trainCount:],
	}

	trainDataset.ItemCount = len(trainDataset.Items)
	valDataset.ItemCount = len(valDataset.Items)

	return trainDataset, valDataset
}

// Filter filters a dataset based on a criteria function
func (d *Dataset) Filter(filterFn func(DatasetItem) bool) *Dataset {
	filtered := &Dataset{
		Path:   d.Path + "_filtered",
		Format: d.Format,
	}

	for _, item := range d.Items {
		if filterFn(item) {
			filtered.Items = append(filtered.Items, item)
		}
	}

	filtered.ItemCount = len(filtered.Items)
	return filtered
}

// Transform applies a transformation to each dataset item
func (d *Dataset) Transform(transformFn func(DatasetItem) DatasetItem) *Dataset {
	transformed := &Dataset{
		Path:   d.Path + "_transformed",
		Format: d.Format,
	}

	for _, item := range d.Items {
		transformed.Items = append(transformed.Items, transformFn(item))
	}

	transformed.ItemCount = len(transformed.Items)
	return transformed
}
