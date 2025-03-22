// Package ui provides user interface components and setup for the RevEnGo application.
// This file contains controllers for managing UI operations and data flow.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"github.com/leog/RevEnGo/internal/models"
	"github.com/leog/RevEnGo/internal/ui/components"
)

// NoteController manages operations related to notes
type NoteController struct {
	noteStore models.NoteStore
	window    fyne.Window
	notepad   fyne.CanvasObject
	sidebar   fyne.CanvasObject

	// Currently loaded note ID (empty if creating a new note)
	currentNoteID string
}

// NewNoteController creates a new controller for note operations
func NewNoteController(noteStore models.NoteStore, window fyne.Window, notepad fyne.CanvasObject, sidebar fyne.CanvasObject) *NoteController {
	return &NoteController{
		noteStore: noteStore,
		window:    window,
		notepad:   notepad,
		sidebar:   sidebar,
	}
}

// CreateNewNote initializes the notepad for creating a new note
func (c *NoteController) CreateNewNote() {
	// Clear the current note ID
	c.currentNoteID = ""

	// Clear the notepad
	// Implementation depends on the notepad's API
}

// SaveCurrentNote saves the current content of the notepad
func (c *NoteController) SaveCurrentNote() error {
	// Extract data from the notepad
	data := components.GetNoteData(c.notepad.(*fyne.Container))

	// Convert to a Note model
	note := components.ConvertToNote(data, c.currentNoteID)

	// Save the note
	err := c.noteStore.SaveNote(note)
	if err != nil {
		dialog.ShowError(err, c.window)
		return err
	}

	// Update current note ID
	c.currentNoteID = note.ID

	// Refresh the sidebar
	// Implementation depends on the sidebar's API

	return nil
}

// LoadNote loads a note into the notepad
func (c *NoteController) LoadNote(noteID string) error {
	// Load the note from storage
	note, err := c.noteStore.GetNote(noteID)
	if err != nil {
		dialog.ShowError(err, c.window)
		return err
	}

	// Convert to NotePadData
	data := components.ConvertFromNote(note)

	// Load data into the notepad
	components.LoadNoteData(c.notepad.(*fyne.Container), data)

	// Update current note ID
	c.currentNoteID = noteID

	return nil
}

// DeleteNote deletes the current note
func (c *NoteController) DeleteNote() error {
	// Confirm deletion
	dialog.ShowConfirm("Delete Note", "Are you sure you want to delete this note?", func(confirmed bool) {
		if confirmed && c.currentNoteID != "" {
			// Delete the note
			err := c.noteStore.DeleteNote(c.currentNoteID)
			if err != nil {
				dialog.ShowError(err, c.window)
				return
			}

			// Clear the notepad
			c.CreateNewNote()

			// Refresh the sidebar
			// Implementation depends on the sidebar's API
		}
	}, c.window)

	return nil
}

// RefreshNoteList updates the sidebar with the current list of notes
func (c *NoteController) RefreshNoteList() error {
	// Implementation depends on the sidebar's API
	return nil
}
