// Package main is the entry point for the RevEnGo application.
// RevEnGo is a note-taking application specifically designed for reverse engineering tasks.
// It provides a clean, intuitive interface for organizing and documenting findings during
// the reverse engineering process.
package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"github.com/leog/RevEnGo/internal/models"
	"github.com/leog/RevEnGo/internal/ui/components"
)

// main is the entry point of the application.
// It initializes the GUI, sets up the application directory structure,
// creates the necessary data stores, and assembles the UI components.
func main() {
	// Create the Fyne application instance
	// This is the root object that manages the application lifecycle
	a := app.New()

	// Set the application theme to dark mode for better visibility
	// This is especially helpful during long reverse engineering sessions
	a.Settings().SetTheme(theme.DarkTheme())

	// Create the main application window with a title
	w := a.NewWindow("RevEnGo - Reverse Engineering Notes")

	// Set the initial window size to be comfortable for note-taking
	// 1200x800 provides enough space for the sidebar and content
	w.Resize(fyne.NewSize(1200, 800))

	// Set up the application's data directories
	// The application will store all data in the user's home directory
	// under the .revengo folder for persistence across sessions
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't access the home directory, the application cannot function
		log.Fatalf("Error getting home directory: %v", err)
	}

	// Create paths for storing application data
	appDir := filepath.Join(homeDir, ".revengo")
	notesDir := filepath.Join(appDir, "notes")       // For storing note files
	projectsDir := filepath.Join(appDir, "projects") // For storing project files

	// Initialize the note storage system
	// This will create the directory if it doesn't exist
	noteStore, err := models.NewFileNoteStore(notesDir)
	if err != nil {
		log.Fatalf("Error initializing note store: %v", err)
	}

	// Initialize the project storage system
	// This will create the directory if it doesn't exist
	projectStore, err := models.NewFileProjectStore(projectsDir)
	if err != nil {
		log.Fatalf("Error initializing project store: %v", err)
	}

	// Create the main UI components
	header := components.NewHeader()   // Top navigation and action buttons
	sidebar := components.NewSidebar() // Navigation and organization panel
	notepad := components.NewNotePad() // Main content editing area

	// Create the main content layout with a horizontal split
	// This divides the screen between the sidebar and notepad
	content := container.NewHSplit(
		sidebar,
		notepad,
	)

	// Set the sidebar to take up 20% of the width
	// This provides enough space for navigation while maximizing the notepad area
	content.Offset = 0.2

	// Assemble the complete UI layout with a border container
	// This places the header at the top and the split content in the center
	mainLayout := container.NewBorder(
		header,  // top component
		nil,     // bottom component (none)
		nil,     // left component (none)
		nil,     // right component (none)
		content, // center component
	)

	// Set the assembled layout as the window's content
	w.SetContent(mainLayout)

	// Keep references to the data stores (will be used for active data management later)
	// Currently marked as unused to prevent build errors
	_ = noteStore
	_ = projectStore

	// Set up a handler for when the window is closed
	// This will be used to save any unsaved data before exiting
	w.SetOnClosed(func() {
		// TODO: Implement saving of unsaved data before closing
		// This will prevent data loss when the application is closed
	})

	// Show the window and start the application's main event loop
	// This call is blocking and will only return when the application exits
	w.ShowAndRun()
}
