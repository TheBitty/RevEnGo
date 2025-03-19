// Package models provides data models and storage functionality for the RevEnGo application.
// This file contains the Project model and its associated storage implementation.
package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Project represents a collection of related notes.
// Projects help organize notes by grouping them under a common theme or purpose,
// such as a specific reverse engineering target or analysis task.
type Project struct {
	// ID is the unique identifier for the project
	// This is typically a timestamp-based string
	ID string `json:"id"`

	// Name is the user-visible title of the project
	Name string `json:"name"`

	// Description provides details about the project's purpose
	Description string `json:"description"`

	// Created is the timestamp when the project was first created
	Created time.Time `json:"created"`

	// Modified is the timestamp when the project was last edited
	Modified time.Time `json:"modified"`
}

// ProjectStore defines the interface for project storage operations.
// This interface abstracts the storage mechanism, allowing different
// implementations (file-based, database, cloud storage, etc.) to be used.
type ProjectStore interface {
	// SaveProject persists a project to storage
	SaveProject(project *Project) error

	// GetProject retrieves a project by its ID
	GetProject(id string) (*Project, error)

	// ListProjects retrieves all projects from storage
	ListProjects() ([]*Project, error)

	// DeleteProject removes a project from storage
	DeleteProject(id string) error
}

// FileProjectStore implements ProjectStore using the local filesystem.
// Projects are stored as individual JSON files in a directory.
type FileProjectStore struct {
	// BasePath is the directory where project files are stored
	BasePath string
}

// NewFileProjectStore creates a new file-based project store.
// It ensures the storage directory exists before returning.
//
// Parameters:
//   - basePath: The directory path where projects will be stored
//
// Returns:
//   - A configured FileProjectStore instance
//   - An error if the directory cannot be created
func NewFileProjectStore(basePath string) (*FileProjectStore, error) {
	// Create the storage directory if it doesn't exist
	// This ensures we can write files immediately
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &FileProjectStore{
		BasePath: basePath,
	}, nil
}

// SaveProject saves a project to the filesystem as a JSON file.
// If the project is new (empty ID), it assigns a new ID and creation timestamp.
// For all projects, the modification timestamp is updated to the current time.
//
// Parameters:
//   - project: The project to save
//
// Returns:
//   - An error if the saving operation fails
func (s *FileProjectStore) SaveProject(project *Project) error {
	// Update the modification time to the current time
	// This happens for both new and existing projects
	project.Modified = time.Now()

	// If it's a new project, generate an ID and set creation time
	if project.ID == "" {
		// Use a timestamp-based ID format (YYYYMMDDhhmmss)
		project.ID = time.Now().Format("20060102150405")
		project.Created = time.Now()
	}

	// Convert the project to a formatted JSON string
	// Use indentation for better human readability
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file named with the project's ID
	filename := filepath.Join(s.BasePath, project.ID+".json")
	return os.WriteFile(filename, data, 0644)
}

// GetProject retrieves a project from the filesystem by its ID.
// It reads the corresponding JSON file and parses it into a Project object.
//
// Parameters:
//   - id: The ID of the project to retrieve
//
// Returns:
//   - The requested project, or nil if not found
//   - An error if the reading or parsing operation fails
func (s *FileProjectStore) GetProject(id string) (*Project, error) {
	// Construct the full path to the project file
	filename := filepath.Join(s.BasePath, id+".json")

	// Read the entire file into memory
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data into a Project object
	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// ListProjects retrieves all projects from the filesystem.
// It finds all JSON files in the storage directory and parses them as projects.
//
// Returns:
//   - A slice containing all projects found in storage
//   - An error if the directory reading operation fails
func (s *FileProjectStore) ListProjects() ([]*Project, error) {
	// Find all JSON files in the storage directory
	pattern := filepath.Join(s.BasePath, "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold all the projects, with initial capacity
	// equal to the number of matching files (to avoid reallocations)
	projects := make([]*Project, 0, len(matches))

	// Read and parse each file, adding valid projects to the result
	for _, match := range matches {
		// Read the file content
		data, err := os.ReadFile(match)
		if err != nil {
			// Skip files that can't be read, rather than failing
			continue
		}

		// Parse the JSON into a Project object
		var project Project
		if err := json.Unmarshal(data, &project); err != nil {
			// Skip files with invalid JSON format
			continue
		}

		// Add the project to our result list
		projects = append(projects, &project)
	}

	return projects, nil
}

// DeleteProject removes a project from the filesystem.
// It deletes the corresponding JSON file.
// Note: This does not delete any notes associated with the project.
// Those would need to be handled separately.
//
// Parameters:
//   - id: The ID of the project to delete
//
// Returns:
//   - An error if the deletion operation fails
func (s *FileProjectStore) DeleteProject(id string) error {
	// Construct the full path to the project file
	filename := filepath.Join(s.BasePath, id+".json")

	// Delete the file from the filesystem
	return os.Remove(filename)
}
