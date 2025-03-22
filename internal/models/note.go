// Package models provides data models and storage functionality for the RevEnGo application.
// This file contains the Note model and its associated storage implementation.
package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Reverse Engineering note types
const (
	RETypeGeneral           = "general"
	RETypeFunctionAnalysis  = "function_analysis"
	RETypeStructureAnalysis = "structure_analysis"
	RETypeProtocolAnalysis  = "protocol_analysis"
	RETypeVulnerability     = "vulnerability"
)

// Note represents a note in the application.
// Each note contains metadata (such as title and tags) and the main content.
// Notes can be associated with projects for organization.
type Note struct {
	// ID is the unique identifier for the note
	// This is typically a timestamp-based string
	ID string `json:"id"`

	// Title is the user-visible name of the note
	Title string `json:"title"`

	// Content is the main text body of the note
	Content string `json:"content"`

	// Tags are keywords for categorization and searching
	Tags []string `json:"tags"`

	// Created is the timestamp when the note was first created
	Created time.Time `json:"created"`

	// Modified is the timestamp when the note was last edited
	Modified time.Time `json:"modified"`

	// ProjectID is an optional reference to a project that contains this note
	// If empty, the note is not associated with any project
	ProjectID string `json:"project_id,omitempty"`

	// RE-specific fields
	BinaryName     string   `json:"binary_name,omitempty"`
	FunctionRefs   []string `json:"function_refs,omitempty"`
	AddressRange   string   `json:"address_range,omitempty"`
	RelatedNotes   []string `json:"related_notes,omitempty"`
	ReverseEngType string   `json:"reverse_eng_type,omitempty"`
}

// NoteStore defines the interface for note storage operations.
// This interface abstracts the storage mechanism, allowing different
// implementations (file-based, database, cloud storage, etc.) to be used.
type NoteStore interface {
	// SaveNote persists a note to storage
	SaveNote(note *Note) error

	// GetNote retrieves a note by its ID
	GetNote(id string) (*Note, error)

	// ListNotes retrieves all notes from storage
	ListNotes() ([]*Note, error)

	// DeleteNote removes a note from storage
	DeleteNote(id string) error
}

// FileNoteStore implements NoteStore using the local filesystem.
// Notes are stored as individual JSON files in a directory.
type FileNoteStore struct {
	// BasePath is the directory where note files are stored
	BasePath string
}

// NewFileNoteStore creates a new file-based note store.
// It ensures the storage directory exists before returning.
//
// Parameters:
//   - basePath: The directory path where notes will be stored
//
// Returns:
//   - A configured FileNoteStore instance
//   - An error if the directory cannot be created
func NewFileNoteStore(basePath string) (*FileNoteStore, error) {
	// Create the storage directory if it doesn't exist,
	// This ensures we can write files immediately
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &FileNoteStore{
		BasePath: basePath,
	}, nil
}

// SaveNote saves a note to the filesystem as a JSON file.
// If the note is new (empty ID), it assigns a new ID and creation timestamp.
// For all notes, the modification timestamp is updated to the current time.
//
// Parameters:
//   - note: The note to save
//
// Returns:
//   - An error if the saving operation fails
func (s *FileNoteStore) SaveNote(note *Note) error {
	// Update the modification time to the current time
	// This happens for both new and existing notes
	note.Modified = time.Now()

	// If it's a new note, generate an ID and set creation time
	if note.ID == "" {
		// Use a timestamp-based ID format (YYYYMMDhhmmss)
		note.ID = time.Now().Format("20060102150405")
		note.Created = time.Now()
	}

	// Convert the note to a formatted JSON string
	// Use indentation for better human readability
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file named with the note's ID
	filename := filepath.Join(s.BasePath, note.ID+".json")
	return os.WriteFile(filename, data, 0644)
}

// GetNote retrieves a note from the filesystem by its ID.
// It reads the corresponding JSON file and parses it into a Note object.
//
// Parameters:
//   - id: The ID of the note to retrieve
//
// Returns:
//   - The requested note, or nil if not found
//   - An error if the reading or parsing operation fails
func (s *FileNoteStore) GetNote(id string) (*Note, error) {
	// Construct the full path to the note file
	filename := filepath.Join(s.BasePath, id+".json")

	// Read the entire file into memory
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data into a Note object
	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// ListNotes retrieves all notes from the filesystem.
// It finds all JSON files in the storage directory and parses them as notes.
//
// Returns:
//   - A slice containing all notes found in storage
//   - An error if the directory reading operation fails
func (s *FileNoteStore) ListNotes() ([]*Note, error) {
	// Find all JSON files in the storage directory
	pattern := filepath.Join(s.BasePath, "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold all the notes, with initial capacity
	// equal to the number of matching files (to avoid reallocations)
	notes := make([]*Note, 0, len(matches))

	// Read and parse each file, adding valid notes to the result
	for _, match := range matches {
		// Read the file content
		data, err := os.ReadFile(match)
		if err != nil {
			// Skip files that can't be read, rather than failing
			continue
		}

		// Parse the JSON into a Note object
		var note Note
		if err := json.Unmarshal(data, &note); err != nil {
			// Skip files with invalid JSON format
			continue
		}

		// Add the note to our result list
		notes = append(notes, &note)
	}

	return notes, nil
}

// DeleteNote removes a note from the filesystem.
// It deletes the corresponding JSON file.
//
// Parameters:
//   - id: The ID of the note to delete
//
// Returns:
//   - An error if the deletion operation fails
func (s *FileNoteStore) DeleteNote(id string) error {
	// Construct the full path to the note file
	filename := filepath.Join(s.BasePath, id+".json")

	// Delete the file from the filesystem
	return os.Remove(filename)
}
