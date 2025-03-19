// Package components provides UI components for the RevEnGo application.
// This package contains reusable UI elements that make up the application interface.
package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewHeader creates a new header component for the application.
// The header provides navigation controls and quick access to common actions.
// It contains:
// - The application title/logo
// - Action buttons for creating, saving, and opening notes
// - A search field for finding notes
// - A settings button for application configuration
//
// Returns a canvas object that can be placed in a container.
func NewHeader() fyne.CanvasObject {
	// Create the application title with bold styling
	// This serves as a visual anchor and reinforces the application identity
	title := widget.NewLabelWithStyle(
		"RevEnGo",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	// Create the "New Note" button with an appropriate icon
	// This button will allow users to create a new note quickly
	newNoteBtn := widget.NewButtonWithIcon("New Note", theme.ContentAddIcon(), func() {
		// Note: Implementation will be added later
		// This will create a new blank note in the editor
	})

	// Create the "Save" button with a save icon
	// This button will save the current note to storage
	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		// Note: Implementation will be added later
		// This will persist the current note to disk
	})

	// Create the "Open" button with a folder icon
	// This button will open a file dialog to select a note to open
	openBtn := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {
		// Note: Implementation will be added later
		// This will show a dialog to open an existing note
	})

	// Create a search entry field for finding notes
	// The search will filter notes based on content or metadata
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search notes...")
	searchEntry.SetMinRowsVisible(1) // Single line height

	// Create a settings button (icon only) for application settings
	// This provides access to application configuration options
	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		// Note: Implementation will be added later
		// This will show the settings dialog
	})

	// Arrange the header elements in containers
	// Left side contains the title and main action buttons
	leftItems := container.NewHBox(title, newNoteBtn, saveBtn, openBtn)

	// Right side contains the search field and settings button
	rightItems := container.NewHBox(searchEntry, settingsBtn)

	// Create the overall header layout using a border container
	// This places elements at the left and right edges with nothing in the center
	header := container.NewBorder(
		nil,        // No top component
		nil,        // No bottom component
		leftItems,  // Left component (title and buttons)
		rightItems, // Right component (search and settings)
		nil,        // No center component
	)

	// Add padding around the header elements for visual comfort
	// This creates space between the header elements and the window edges
	return container.NewPadded(header)
}
