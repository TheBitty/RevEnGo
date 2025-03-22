// Package components provides UI components for the RevEnGo application.
// This file contains the notepad component used for creating and editing notes.
package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/leog/RevEnGo/internal/models"
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

	// RE-specific fields
	BinaryName     string
	FunctionRefs   []string
	AddressRange   string
	RelatedNotes   []string
	ReverseEngType string
}

// NewNotePad creates a new notepad component for editing and viewing notes.
// The notepad provides:
// - A title field for naming the note
// - A large content area for the main note text
// - A tags field for categorization
// - RE-specific fields for specialized analysis
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

	// Create RE-specific fields
	// Note type selector
	noteTypeSelect := widget.NewSelect([]string{
		models.RETypeGeneral,
		models.RETypeFunctionAnalysis,
		models.RETypeStructureAnalysis,
		models.RETypeProtocolAnalysis,
		models.RETypeVulnerability,
	}, nil)
	noteTypeSelect.SetSelected(models.RETypeGeneral)

	// Binary name entry
	binaryNameEntry := widget.NewEntry()
	binaryNameEntry.SetPlaceHolder("Binary Name (optional)")

	// Address range entry
	addressRangeEntry := widget.NewEntry()
	addressRangeEntry.SetPlaceHolder("Address Range (e.g., 0x1000-0x2000)")

	// Function references entry
	functionRefsEntry := widget.NewMultiLineEntry()
	functionRefsEntry.SetPlaceHolder("Function references (one per line)")
	functionRefsEntry.SetMinRowsVisible(3)

	// Arrange the tags label and entry field in a horizontal layout
	tagsContainer := container.NewBorder(nil, nil, tagsLabel, nil, tagsEntry)

	// Create container for RE-specific fields
	reFieldsContainer := container.NewVBox(
		widget.NewLabel("Note Type:"),
		noteTypeSelect,
		widget.NewLabel("Binary Name:"),
		binaryNameEntry,
		widget.NewLabel("Address Range:"),
		addressRangeEntry,
		widget.NewLabel("Function References:"),
		functionRefsEntry,
	)

	// Create tabs for regular note fields and RE-specific fields
	tabs := container.NewAppTabs(
		container.NewTabItem("Basic Info", container.NewVBox(
			titleEntry,
			tagsContainer,
		)),
		container.NewTabItem("RE Details", reFieldsContainer),
	)

	// Create the overall notepad layout
	// This places the title at the top, content in the center, and tags at the bottom
	noteContainer := container.NewBorder(
		tabs,         // Top component (tabs)
		nil,          // No bottom component
		nil,          // No left component
		nil,          // No right component
		contentEntry, // Content in the center (the largest area)
	)

	// Add padding around the notepad for visual comfort
	// This creates space between the notepad elements and the container edges
	return container.NewPadded(noteContainer)
}

// LoadNoteData loads data into the notepad component.
// This function populates the notepad with existing note data
// when a user selects a note to view or edit.
//
// Parameters:
//   - notepad: The notepad container to load data into
//   - data: The note data to load
func LoadNoteData(notepad *fyne.Container, data NotePadData) {
	// Implementation to be completed
}

// GetNoteData retrieves data from the notepad component.
// This function extracts the current note data from the UI
// when a user wants to save a note.
//
// Parameters:
//   - notepad: The notepad container to extract data from
//
// Returns:
//   - The note data extracted from the notepad
func GetNoteData(notepad *fyne.Container) NotePadData {
	// Implementation to be completed
	return NotePadData{}
}

// ConvertToNote converts NotePadData to a models.Note.
// This function is used when saving the current UI data to storage.
func ConvertToNote(data NotePadData, existingID string) *models.Note {
	note := &models.Note{
		ID:             existingID,
		Title:          data.Title,
		Content:        data.Content,
		Tags:           data.Tags,
		BinaryName:     data.BinaryName,
		FunctionRefs:   data.FunctionRefs,
		AddressRange:   data.AddressRange,
		RelatedNotes:   data.RelatedNotes,
		ReverseEngType: data.ReverseEngType,
	}
	return note
}

// ConvertFromNote converts a models.Note to NotePadData.
// This function is used when loading stored data into the UI.
func ConvertFromNote(note *models.Note) NotePadData {
	return NotePadData{
		Title:          note.Title,
		Content:        note.Content,
		Tags:           note.Tags,
		BinaryName:     note.BinaryName,
		FunctionRefs:   note.FunctionRefs,
		AddressRange:   note.AddressRange,
		RelatedNotes:   note.RelatedNotes,
		ReverseEngType: note.ReverseEngType,
	}
}
