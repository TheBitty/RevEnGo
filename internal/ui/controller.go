// Package ui provides user interface components and setup for the RevEnGo application.
// This file contains controllers for managing UI operations and data flow.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

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
	components.ClearNotepad(c.notepad.(*fyne.Container))
}

// SaveCurrentNote saves the current content of the notepad
func (c *NoteController) SaveCurrentNote() error {
	// Extract data from the notepad
	data := components.GetNoteData(c.notepad.(*fyne.Container))

	// Validate data
	if data.Title == "" {
		dialog.ShowInformation("Missing Information", "Please provide a title for your note.", c.window)
		return nil
	}

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
	c.RefreshNoteList()

	// Show success message
	dialog.ShowInformation("Note Saved", "Your note has been saved successfully.", c.window)

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
	if c.currentNoteID == "" {
		dialog.ShowInformation("No Note Selected", "Please select a note to delete.", c.window)
		return nil
	}

	// Confirm deletion
	dialog.ShowConfirm("Delete Note", "Are you sure you want to delete this note?", func(confirmed bool) {
		if confirmed {
			// Delete the note
			err := c.noteStore.DeleteNote(c.currentNoteID)
			if err != nil {
				dialog.ShowError(err, c.window)
				return
			}

			// Clear the notepad
			c.CreateNewNote()

			// Refresh the sidebar
			c.RefreshNoteList()

			// Show success message
			dialog.ShowInformation("Note Deleted", "The note has been deleted successfully.", c.window)
		}
	}, c.window)

	return nil
}

// RefreshNoteList updates the sidebar with the current list of notes
func (c *NoteController) RefreshNoteList() error {
	// Get all notes
	notes, err := c.noteStore.ListNotes()
	if err != nil {
		dialog.ShowError(err, c.window)
		return err
	}

	var content fyne.CanvasObject

	if len(notes) == 0 {
		// No notes yet, show message
		content = container.NewVBox(
			widget.NewLabel("Your Notes"),
			widget.NewLabel("No notes yet. Create one using the toolbar!"),
		)
	} else {
		// Create a list of items for the sidebar
		notesList := widget.NewList(
			func() int {
				return len(notes)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				obj.(*widget.Label).SetText(notes[id].Title)
			},
		)

		// Set up on-selected handler
		notesList.OnSelected = func(id widget.ListItemID) {
			if id < len(notes) {
				c.LoadNote(notes[id].ID)
			}
		}

		// Wrap in a container with header
		content = container.NewBorder(
			widget.NewLabel("Your Notes"),
			nil,
			nil,
			nil,
			notesList,
		)
	}

	// Update the sidebar using the component's function
	components.UpdateNotesList(c.sidebar.(*fyne.Container), content)

	return nil
}
