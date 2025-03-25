package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"

	"github.com/leog/RevEnGo/internal/models"
	"github.com/leog/RevEnGo/internal/ui"
	"github.com/leog/RevEnGo/internal/ui/theme"
)

func main() {
	// This is the root object that manages the application lifecycle
	a := app.New()

	// Set up custom RevEnGo theme
	a.Settings().SetTheme(theme.New())

	// Create the main application window with a title
	w := a.NewWindow("RevEnGo")

	homeDir, err := os.UserHomeDir() // The application will store all data in the user's home directory
	if err != nil {
		// If we can't access the home directory, the application cannot function
		log.Fatalf("Error getting home directory: %v", err)
	}

	// Paths for storing application data
	appDir := filepath.Join(homeDir, ".revengo")
	notesDir := filepath.Join(appDir, "notes")       // For storing note files
	projectsDir := filepath.Join(appDir, "projects") // For storing project files

	// Ensure program flow directory exists (used for storing program flow diagrams)
	programFlowDir := filepath.Join(homeDir, "Program Flow")
	if err := os.MkdirAll(programFlowDir, 0755); err != nil {
		log.Printf("Warning: Failed to create program flow directory: %v", err)
	}

	// Initialize the note storage system
	noteStore, err := models.NewFileNoteStore(notesDir)
	if err != nil {
		log.Fatalf("Error initializing note store: %v", err)
	}

	// Initialize the project storage system
	projectStore, err := models.NewFileProjectStore(projectsDir)
	if err != nil {
		log.Fatalf("Error initializing project store: %v", err)
	}

	// Create config for UI setup
	config := ui.AppConfig{
		NoteStore:    noteStore,
		ProjectStore: projectStore,
	}

	// Set up the main window with the configuration
	ui.SetupMainWindow(w, config)

	// Start the application
	w.ShowAndRun()
}
