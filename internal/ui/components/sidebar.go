// Package components provide UI components for the RevEnGo application.
// This file contains the sidebar component used for navigation and organization.
package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/leog/RevEnGo/internal/models"
	"github.com/leog/RevEnGo/internal/ui/widgets"
)

// Sidebar color constants
var (
	sidebarBgColor = color.NRGBA{R: 12, G: 17, B: 27, A: 255}   // Dark blue-black
	sidebarAccent  = color.NRGBA{R: 0, G: 174, B: 239, A: 255}  // Cyber blue accent
	highlightColor = color.NRGBA{R: 30, G: 60, B: 110, A: 255}  // Selection highlight
	terminalGreen  = color.NRGBA{R: 35, G: 209, B: 139, A: 255} // Terminal green
)

// SidebarSection represents a section in the sidebar navigation.
// Each section contains a title, icon, and a list of items that can be selected.
// Sections help organize related items into collapsible groups.
type SidebarSection struct {
	// The Title is the display name of the section
	Title string

	// Icon is the visual representation of the section
	Icon fyne.Resource

	// Items is the list of selectable entries in this section
	Items []string
}

// NewSidebar creates a new sidebar component for navigation.
// The sidebar provides a hierarchical structure for organizing and accessing notes.
// It displays sections for:
// - Notes for accessing all notes
// - Recent notes for quick access to recently edited items
// - Projects for organizing related notes
// - Tags for filtering notes by keywords
//
// Returns a canvas object that can be placed in a container.
func NewSidebar() fyne.CanvasObject {
	// Create background panel
	background := canvas.NewRectangle(sidebarBgColor)

	// Create the sidebar title with digital styling
	sidebarTitle := widgets.DigitalText("OPERATIONS", 16, sidebarAccent)

	// Create binary decoration
	binaryDecoration := canvas.NewText("01001100 01001111 01000111", color.NRGBA{R: 40, G: 80, B: 120, A: 120})
	binaryDecoration.TextSize = 9
	binaryDecoration.TextStyle = fyne.TextStyle{Monospace: true}

	// Create a separator with circuit-inspired styling
	circuitLine := canvas.NewLine(sidebarAccent)
	circuitLine.StrokeWidth = 1

	// Create hexagonal navigation buttons with icons
	notesButton := widgets.HexagonalButton(theme.DocumentIcon(), nil)
	projectsButton := widgets.HexagonalButton(theme.FolderIcon(), nil)
	analysisButton := widgets.HexagonalButton(theme.ViewRestoreIcon(), nil)
	settingsButton := widgets.HexagonalButton(theme.SettingsIcon(), nil)

	// Create a toolbar with the hexagonal buttons
	hexButtonsContainer := container.NewHBox(
		notesButton,
		projectsButton,
		analysisButton,
		settingsButton,
	)

	// Create a container for the notes list (placeholder)
	notesListContainer := container.NewVBox(
		widget.NewLabelWithStyle("NOTES", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true}),
		createNoteTypeIndicator("function_analysis", "Stack Buffer Analysis"),
		createNoteTypeIndicator("vulnerability", "Heap Overflow CVE-2023-1234"),
		createNoteTypeIndicator("structure_analysis", "PE Header Structure"),
		createNoteTypeIndicator("general", "Project Overview"),
	)

	// Create header section with title and decoration
	headerSection := container.NewVBox(
		sidebarTitle,
		binaryDecoration,
		circuitLine,
		container.NewPadded(hexButtonsContainer),
	)

	// Combine all elements into a vertical layout
	sidebarContent := container.NewBorder(
		headerSection,
		nil,
		nil,
		nil,
		container.NewPadded(notesListContainer),
	)

	// Stack the background and content
	return container.NewStack(
		background,
		sidebarContent,
	)
}

// UpdateNotesList updates the notes list in the sidebar
// This function is intended to be called by the controller when the list of notes changes
func UpdateNotesList(sidebar *fyne.Container, notesList fyne.CanvasObject) {
	// The sidebar is now a stack with background and content layers
	// We need to update the content layer's center component
	contentContainer := sidebar.Objects[1].(*fyne.Container)

	// The content container is a border layout, replace the center component (notes list)
	contentContainer.Objects[0] = container.NewBorder(
		contentContainer.Objects[0].(*fyne.Container), // Keep the existing header
		nil,
		nil,
		nil,
		container.NewPadded(notesList), // Update the notes list
	)

	contentContainer.Refresh()
}

// createSidebarTree creates a tree widget from the sidebar sections.
// It converts the sections and their items into a hierarchical tree structure
// that can be displayed in the sidebar.
//
// Parameters:
//   - sections: The sections to include in the tree
//
// Returns:
//   - A tree widget configured with the provided sections
func createSidebarTree(sections []SidebarSection) *widget.Tree {
	// Create a new tree widget with the necessary callback functions
	tree := widget.NewTree(
		// This function defines the child IDs for a given node ID
		// It builds the tree structure by returning children for each node
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			// For the root level, return the section titles as children
			if id == "" {
				ids := make([]widget.TreeNodeID, len(sections))
				for i, section := range sections {
					ids[i] = section.Title
				}
				return ids
			}

			// For section nodes, return their items as children
			// Items are prefixed with the section name for uniqueness
			for _, section := range sections {
				if id == section.Title {
					ids := make([]widget.TreeNodeID, len(section.Items))
					for i, item := range section.Items {
						ids[i] = section.Title + "." + item
					}
					return ids
				}
			}

			// Return empty slice for leaf nodes
			return []widget.TreeNodeID{}
		},

		// This function determines if a node is a branch (can have children)
		// In our case, only section titles are branches
		func(id widget.TreeNodeID) bool {
			// If the ID contains a dot, it's an item (leaf) not a section (branch)
			for _, r := range id {
				if r == '.' {
					return false
				}
			}
			return true
		},

		// This function creates the template for each type of node
		// It defines how branch nodes and leaf nodes should appear
		func(branch bool) fyne.CanvasObject {
			if branch {
				// For branch nodes (sections), show a folder icon
				icon := widget.NewIcon(theme.FolderIcon())
				return container.NewHBox(icon, widget.NewLabel("Section"))
			}
			// For leaf nodes (items), show a document icon
			icon := widget.NewIcon(theme.DocumentIcon())
			return container.NewHBox(icon, widget.NewLabel("Item"))
		},

		// This function updates nodes with their specific content
		// It sets the correct label text and icon for each node
		func(id widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			nodeContainer := obj.(*fyne.Container)
			label := nodeContainer.Objects[1].(*widget.Label)
			icon := nodeContainer.Objects[0].(*widget.Icon)

			if branch {
				// For branch nodes, use the section title and icon
				for _, section := range sections {
					if id == section.Title {
						label.SetText(section.Title)
						icon.SetResource(section.Icon)
						return
					}
				}

			} else {
				// For leaf nodes, extract the item name from the ID
				// (everything after the dot)
				for i := len(id) - 1; i >= 0; i-- {
					if id[i] == '.' {
						label.SetText(id[i+1:])
						icon.SetResource(theme.DocumentIcon())
						return
					}
				}
			}
		},
	)

	// Expand all section branches by default for better visibility
	// This shows all items under each section when the sidebar first loads
	for _, section := range sections {
		tree.OpenBranch(section.Title)
	}

	return tree
}

// createNoteTypeIndicator creates a list item with an indicator showing the type of note
func createNoteTypeIndicator(noteType string, title string) fyne.CanvasObject {
	var indicatorColor color.Color
	var iconRes fyne.Resource

	// Choose color and icon based on note type
	switch noteType {
	case models.RETypeFunctionAnalysis:
		indicatorColor = color.NRGBA{R: 0, G: 180, B: 255, A: 255} // Blue
		iconRes = theme.DocumentIcon()
	case models.RETypeVulnerability:
		indicatorColor = color.NRGBA{R: 255, G: 70, B: 70, A: 255} // Red
		iconRes = theme.WarningIcon()
	case models.RETypeStructureAnalysis:
		indicatorColor = color.NRGBA{R: 180, G: 120, B: 255, A: 255} // Purple
		iconRes = theme.StorageIcon()
	case models.RETypeProtocolAnalysis:
		indicatorColor = color.NRGBA{R: 255, G: 180, B: 0, A: 255} // Amber
		iconRes = theme.MailComposeIcon()
	default:
		indicatorColor = color.NRGBA{R: 120, G: 120, B: 120, A: 255} // Gray
		iconRes = theme.DocumentIcon()
	}

	// Create an icon with the appropriate color
	icon := widget.NewIcon(iconRes)

	// Create a color indicator
	indicator := canvas.NewRectangle(indicatorColor)
	indicator.SetMinSize(fyne.NewSize(4, 20))

	// Create the title label with monospaced font
	label := widget.NewLabel(title)
	label.TextStyle = fyne.TextStyle{Monospace: true}

	// Create the item container
	itemContent := container.NewBorder(
		nil,
		nil,
		container.NewHBox(indicator, icon),
		nil,
		label,
	)

	// Create the hoverable container
	hoverRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	item := container.NewStack(hoverRect, itemContent)

	// Make the item visually distinct
	hoverRect.FillColor = color.NRGBA{R: 15, G: 30, B: 55, A: 100}

	return container.NewPadded(item)
}
