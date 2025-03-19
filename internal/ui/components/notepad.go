// Package components provides UI components for the RevEnGo application.
// This file contains the notepad component used for creating and editing notes.
package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NotePadData represents the data for a note in the application.
// This structure encapsulates all the necessary information for a single note,
// including its title, content, and associated tags.
type NotePadData struct {
	// Title is the name/header of the note
	Title string

	// Content is the main text body of the note
	Content string

	// Tags is a list of keywords used for categorization and searching
	Tags []string
}

// NewNotePad creates a new notepad component for editing and viewing notes.
// The notepad provides:
// - A title field for naming the note
// - A large content area for the main note text
// - A tags field for categorization
//
// Returns a canvas object that can be placed in a container.
func NewNotePad() fyne.CanvasObject {
	// Create the title entry field for the note's name
	// This field is prominent at the top of the notepad
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Note Title")

	// Create the main content entry field for the note's body
	// This is a multi-line text area where the primary note content is written
	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetPlaceHolder("Write your analysis notes here...")
	contentEntry.SetMinRowsVisible(20) // Set a comfortable height for writing

	// Create the tags entry field for categorizing the note
	// Tags help with organization and searching for notes later
	tagsEntry := widget.NewEntry()
	tagsEntry.SetPlaceHolder("Tags (comma separated)")

	// Create a label for the tags field to clearly identify its purpose
	tagsLabel := widget.NewLabel("Tags:")

	// Arrange the tags label and entry field in a horizontal layout
	tagsContainer := container.NewBorder(nil, nil, tagsLabel, nil, tagsEntry)

	// Create the overall notepad layout
	// This places the title at the top, content in the center, and tags at the bottom
	noteContainer := container.NewBorder(
		container.NewVBox(
			titleEntry,    // Title at the top
			tagsContainer, // Tags below the title
		),
		nil,          // No bottom component
		nil,          // No left component
		nil,          // No right component
		contentEntry, // Content in the center (largest area)
	)

	// Add padding around the notepad for visual comfort
	// This creates space between the notepad elements and the container edges
	return container.NewPadded(noteContainer)
}

// LoadNoteData loads data into the notepad component.
// This function will be used to populate the notepad with existing note data
// when a user selects a note to view or edit.
//
// Parameters:
//   - notepad: The notepad container to load data into
//   - data: The note data to load
func LoadNoteData(notepad *fyne.Container, data NotePadData) {
	// Note: Implementation will be added later
	// This will set the title, content, and tags fields with the provided data
}

// GetNoteData retrieves data from the notepad component.
// This function will be used to extract the current note data from the UI
// when a user wants to save a note.
//
// Parameters:
//   - notepad: The notepad container to extract data from
//
// Returns:
//   - The note data extracted from the notepad
func GetNoteData(notepad *fyne.Container) NotePadData {
	// Note: Implementation will be added later
	// This will retrieve the title, content, and tags from the UI fields
	return NotePadData{}
}
