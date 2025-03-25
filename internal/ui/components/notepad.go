// Package components provides UI components for the RevEnGo application.
// This file contains the notepad component used for creating and editing notes.
package components

import (
	"strings"

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

// NotePadComponents groups all the interactive components of the notepad
// This makes it easier to access them for setting and getting data
type NotePadComponents struct {
	TitleEntry        *widget.Entry
	ContentEntry      *widget.Entry
	TagsEntry         *widget.Entry
	NoteTypeSelect    *widget.Select
	BinaryNameEntry   *widget.Entry
	AddressRangeEntry *widget.Entry
	FunctionRefsEntry *widget.Entry
	Tabs              *container.AppTabs
}

// Global reference to the current notepad components
// This approach avoids the need to store components in the container
var currentComponents *NotePadComponents

// NewNotePad creates a new notepad component for editing and viewing notes.
// The notepad provides:
// - A title field for naming the note
// - A large content area for the main note text
// - A tags field for categorization
// - RE-specific fields for specialized analysis
//
// Returns a canvas object that can be placed in a container.
func NewNotePad() fyne.CanvasObject {
	components := &NotePadComponents{}
	currentComponents = components // Store the reference globally

	// Create the title entry field for the note's name
	components.TitleEntry = widget.NewEntry()
	components.TitleEntry.SetPlaceHolder("Note Title")

	// Create the main content entry field for the note's body
	components.ContentEntry = widget.NewMultiLineEntry()
	components.ContentEntry.SetPlaceHolder("Write your analysis notes here...")
	components.ContentEntry.SetMinRowsVisible(20) // Set a comfortable height for writing

	// Create the tags entry field for categorizing the note
	components.TagsEntry = widget.NewEntry()
	components.TagsEntry.SetPlaceHolder("Tags (comma separated)")

	// Create a label for the tags field
	tagsLabel := widget.NewLabel("Tags:")

	// Create RE-specific fields
	// Note type selector
	components.NoteTypeSelect = widget.NewSelect([]string{
		models.RETypeGeneral,
		models.RETypeFunctionAnalysis,
		models.RETypeStructureAnalysis,
		models.RETypeProtocolAnalysis,
		models.RETypeVulnerability,
	}, nil)
	components.NoteTypeSelect.SetSelected(models.RETypeGeneral)

	// Binary name entry
	components.BinaryNameEntry = widget.NewEntry()
	components.BinaryNameEntry.SetPlaceHolder("Binary Name (optional)")

	// Address range entry
	components.AddressRangeEntry = widget.NewEntry()
	components.AddressRangeEntry.SetPlaceHolder("Address Range (e.g., 0x1000-0x2000)")

	// Function references entry
	components.FunctionRefsEntry = widget.NewMultiLineEntry()
	components.FunctionRefsEntry.SetPlaceHolder("Function references (one per line)")
	components.FunctionRefsEntry.SetMinRowsVisible(3)

	// Arrange the tags label and entry field in a horizontal layout
	tagsContainer := container.NewBorder(nil, nil, tagsLabel, nil, components.TagsEntry)

	// Create container for RE-specific fields
	reFieldsContainer := container.NewVBox(
		widget.NewLabel("Note Type:"),
		components.NoteTypeSelect,
		widget.NewLabel("Binary Name:"),
		components.BinaryNameEntry,
		widget.NewLabel("Address Range:"),
		components.AddressRangeEntry,
		widget.NewLabel("Function References:"),
		components.FunctionRefsEntry,
	)

	// Create tabs for regular note fields and RE-specific fields
	components.Tabs = container.NewAppTabs(
		container.NewTabItem("Basic Info", container.NewVBox(
			components.TitleEntry,
			tagsContainer,
		)),
		container.NewTabItem("RE Details", reFieldsContainer),
	)

	// Create the overall notepad layout
	noteContainer := container.NewBorder(
		components.Tabs,         // Top component (tabs)
		nil,                     // No bottom component
		nil,                     // No left component
		nil,                     // No right component
		components.ContentEntry, // Content in the center (the largest area)
	)

	// Add padding around the notepad for visual comfort
	return container.NewPadded(noteContainer)
}

// getComponents returns the current components
func getComponents() *NotePadComponents {
	return currentComponents
}

// LoadNoteData loads data into the notepad component.
// This function populates the notepad with existing note data
// when a user selects a note to view or edit.
//
// Parameters:
//   - notepad: The notepad container to load data into
//   - data: The note data to load
func LoadNoteData(notepad *fyne.Container, data NotePadData) {
	components := getComponents()

	// Set basic note data
	components.TitleEntry.SetText(data.Title)
	components.ContentEntry.SetText(data.Content)
	components.TagsEntry.SetText(strings.Join(data.Tags, ", "))

	// Set RE-specific data
	components.NoteTypeSelect.SetSelected(data.ReverseEngType)
	components.BinaryNameEntry.SetText(data.BinaryName)
	components.AddressRangeEntry.SetText(data.AddressRange)
	components.FunctionRefsEntry.SetText(strings.Join(data.FunctionRefs, "\n"))
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
	components := getComponents()

	// Extract tags from comma-separated list
	tags := []string{}
	if components.TagsEntry.Text != "" {
		for _, tag := range strings.Split(components.TagsEntry.Text, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	// Extract function references from newline-separated list
	functionRefs := []string{}
	if components.FunctionRefsEntry.Text != "" {
		for _, ref := range strings.Split(components.FunctionRefsEntry.Text, "\n") {
			if trimmed := strings.TrimSpace(ref); trimmed != "" {
				functionRefs = append(functionRefs, trimmed)
			}
		}
	}

	// Compile the data
	return NotePadData{
		Title:          components.TitleEntry.Text,
		Content:        components.ContentEntry.Text,
		Tags:           tags,
		BinaryName:     components.BinaryNameEntry.Text,
		FunctionRefs:   functionRefs,
		AddressRange:   components.AddressRangeEntry.Text,
		ReverseEngType: components.NoteTypeSelect.Selected,
	}
}

// ClearNotepad resets all fields in the notepad
func ClearNotepad(notepad *fyne.Container) {
	components := getComponents()

	// Clear basic note data
	components.TitleEntry.SetText("")
	components.ContentEntry.SetText("")
	components.TagsEntry.SetText("")

	// Clear RE-specific data
	components.NoteTypeSelect.SetSelected(models.RETypeGeneral)
	components.BinaryNameEntry.SetText("")
	components.AddressRangeEntry.SetText("")
	components.FunctionRefsEntry.SetText("")

	// Reset to first tab
	components.Tabs.SelectIndex(0)
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
