// Package components provides UI components for the RevEnGo application.
// This file contains the notepad component used for creating and editing notes.
package components

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/leog/RevEnGo/internal/models"
	"github.com/leog/RevEnGo/internal/ui/widgets"
)

// Color constants for the notepad
var (
	terminalBgColor    = color.NRGBA{R: 8, G: 14, B: 21, A: 255}     // Dark terminal background
	terminalTextColor  = color.NRGBA{R: 180, G: 255, B: 180, A: 255} // Terminal green text
	codeBlockBgColor   = color.NRGBA{R: 15, G: 25, B: 35, A: 255}    // Slightly lighter for code blocks
	tabActiveBgColor   = color.NRGBA{R: 30, G: 50, B: 70, A: 255}    // Active tab background
	tabInactiveBgColor = color.NRGBA{R: 15, G: 30, B: 50, A: 255}    // Inactive tab background
	accentBlue         = color.NRGBA{R: 0, G: 174, B: 239, A: 255}   // Cyber blue accent
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

	// Create the background
	background := canvas.NewRectangle(terminalBgColor)

	// Create the title entry field with terminal styling
	components.TitleEntry = widget.NewEntry()
	components.TitleEntry.SetPlaceHolder("Note Title")
	components.TitleEntry.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	// Create the main content entry field with terminal styling
	components.ContentEntry = widgets.TerminalEntry()
	components.ContentEntry.SetPlaceHolder("Write your analysis notes here...")
	components.ContentEntry.SetMinRowsVisible(20)

	// Create the tags entry field with terminal styling
	components.TagsEntry = widget.NewEntry()
	components.TagsEntry.SetPlaceHolder("Tags (comma separated)")
	components.TagsEntry.TextStyle = fyne.TextStyle{Monospace: true}

	// Create a label for the tags field with distinctive styling
	tagsLabel := canvas.NewText("TAGS:", color.NRGBA{R: 0, G: 200, B: 170, A: 255})
	tagsLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	tagsLabel.TextSize = 12

	// Create RE-specific fields with terminal styling
	// Note type selector with distinctive styling
	components.NoteTypeSelect = widget.NewSelect([]string{
		models.RETypeGeneral,
		models.RETypeFunctionAnalysis,
		models.RETypeStructureAnalysis,
		models.RETypeProtocolAnalysis,
		models.RETypeVulnerability,
	}, nil)
	components.NoteTypeSelect.SetSelected(models.RETypeGeneral)

	// Binary name entry with terminal styling
	components.BinaryNameEntry = widget.NewEntry()
	components.BinaryNameEntry.SetPlaceHolder("Binary Name (optional)")
	components.BinaryNameEntry.TextStyle = fyne.TextStyle{Monospace: true}

	// Address range entry with terminal styling
	components.AddressRangeEntry = widget.NewEntry()
	components.AddressRangeEntry.SetPlaceHolder("Address Range (e.g., 0x1000-0x2000)")
	components.AddressRangeEntry.TextStyle = fyne.TextStyle{Monospace: true}

	// Function references entry with terminal styling
	components.FunctionRefsEntry = widget.NewMultiLineEntry()
	components.FunctionRefsEntry.SetPlaceHolder("Function references (one per line)")
	components.FunctionRefsEntry.SetMinRowsVisible(3)
	components.FunctionRefsEntry.TextStyle = fyne.TextStyle{Monospace: true}

	// Create title container with prompt-like styling
	titlePrompt := canvas.NewText(">> ", accentBlue)
	titlePrompt.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	titlePrompt.TextSize = 16

	titleContainer := container.NewBorder(
		nil, nil, titlePrompt, nil, components.TitleEntry)

	// Arrange the tags label and entry field in a horizontal layout
	tagsContainer := container.NewBorder(nil, nil, tagsLabel, nil, components.TagsEntry)

	// Create styled labels for RE fields
	typeLabel := createTerminalLabel("TYPE:")
	binaryLabel := createTerminalLabel("BINARY:")
	addressLabel := createTerminalLabel("ADDR_RANGE:")
	funcRefsLabel := createTerminalLabel("XREFS:")

	// Create container for RE-specific fields with terminal styling
	reFieldsContainer := container.NewVBox(
		container.NewBorder(nil, nil, typeLabel, nil, components.NoteTypeSelect),
		container.NewBorder(nil, nil, binaryLabel, nil, components.BinaryNameEntry),
		container.NewBorder(nil, nil, addressLabel, nil, components.AddressRangeEntry),
		funcRefsLabel,
		components.FunctionRefsEntry,
	)

	// Create a code block background for the RE fields
	reBackground := canvas.NewRectangle(codeBlockBgColor)
	reContainer := container.NewStack(
		reBackground,
		container.NewPadded(reFieldsContainer),
	)

	// Create tabs for regular note fields and RE-specific fields with hex addresses
	components.Tabs = container.NewAppTabs(
		widgets.HexTabItem("01", container.NewVBox(
			titleContainer,
			tagsContainer,
		)),
		widgets.HexTabItem("02", reContainer),
	)

	// Style the tabs to look like memory addresses
	components.Tabs.OnSelected = func(tab *container.TabItem) {
		// We could add more custom tab styling logic here
	}

	// Add decorative elements to make it look like a terminal
	// Create a terminal-style prompt for the content area
	contentPrompt := canvas.NewText("$>", accentBlue)
	contentPrompt.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	contentPrompt.TextSize = 14

	// Add hex address indicators to simulate memory view
	addrIndicator := createHexAddressLabel()

	// Create the content area with decorative elements
	contentContainer := container.NewBorder(
		container.NewHBox(contentPrompt, addrIndicator),
		nil,
		nil,
		nil,
		components.ContentEntry,
	)

	// Create the overall notepad layout
	noteContainer := container.NewBorder(
		components.Tabs,  // Top component (tabs)
		nil,              // No bottom component
		nil,              // No left component
		nil,              // No right component
		contentContainer, // Content in the center
	)

	// Stack the background and content
	return container.NewStack(
		background,
		container.NewPadded(noteContainer),
	)
}

// createTerminalLabel creates a terminal-styled label
func createTerminalLabel(text string) *canvas.Text {
	label := canvas.NewText(text, accentBlue)
	label.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	label.TextSize = 12
	return label
}

// createHexAddressLabel creates a label with hexadecimal address styling
func createHexAddressLabel() *canvas.Text {
	label := canvas.NewText("0x00c0ffee:", terminalTextColor)
	label.TextStyle = fyne.TextStyle{Monospace: true}
	label.TextSize = 12
	return label
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
