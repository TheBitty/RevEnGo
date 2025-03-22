// Package ui provides user interface components and setup for the RevEnGo application.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/leog/RevEnGo/internal/models"
	"github.com/leog/RevEnGo/internal/ui/components"
)

// AppConfig holds the application configuration and dependencies
type AppConfig struct {
	NoteStore    models.NoteStore
	ProjectStore models.ProjectStore
}

// SetupMainWindow configures the main application window and its components
func SetupMainWindow(w fyne.Window, config AppConfig) {
	// Set window size
	w.Resize(fyne.NewSize(1200, 800))

	// Create the main UI components
	header := components.NewHeader()
	sidebar := components.NewSidebar()
	notepad := components.NewNotePad()

	// Create the content layout
	content := container.NewHSplit(
		sidebar,
		notepad,
	)
	content.Offset = 0.2

	// Create the main layout
	mainLayout := container.NewBorder(
		header,  // top component
		nil,     // bottom component (none)
		nil,     // left component (none)
		nil,     // right component (none)
		content, // center component
	)

	// Create note controller
	noteController := NewNoteController(config.NoteStore, w, notepad, sidebar)

	// Set up toolbar actions
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			noteController.CreateNewNote()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			noteController.SaveCurrentNote()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			noteController.DeleteNote()
		}),
	)

	// Add toolbar to the header
	headerContainer := container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		header,
	)

	// Update the main layout with the new header
	mainLayout = container.NewBorder(
		headerContainer,
		nil,
		nil,
		nil,
		content,
	)

	// Set the window content
	w.SetContent(mainLayout)

	// Set up window close handler
	w.SetOnClosed(func() {
		// TODO: Implement saving of unsaved data before closing
	})

	// Load initial note list
	noteController.RefreshNoteList()
}

// SetupAppTheme configures the application theme
func SetupAppTheme(a fyne.App) {
	a.Settings().SetTheme(theme.DarkTheme())
}
