// Package components provide UI components for the RevEnGo application.
// This file contains the sidebar component used for navigation and organization.
package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
// - Recent notes for quick access to recently edited items
// - Projects for organizing related notes
// - Tags for filtering notes by keywords
//
// Returns a canvas object that can be placed in a container.
func NewSidebar() fyne.CanvasObject {
	// Create the "Recent Notes" section
	// This section displays recently edited notes for quick access
	recentNotes := SidebarSection{
		Title: "Recent Notes",
		Icon:  theme.DocumentIcon(),
		Items: []string{"Note 1", "Note 2", "Note 3"}, // Example items
	}

	// Create the "Projects" section
	// This section organizes notes by project/category
	projects := SidebarSection{
		Title: "Projects",
		Icon:  theme.FolderIcon(),
		Items: []string{"Project A", "Project B"}, // Example projects
	}

	// Create the "Tags" section
	// This section allows filtering notes by tags
	tags := SidebarSection{
		Title: "Tags",
		Icon:  theme.InfoIcon(),
		Items: []string{"Binary", "Assembly", "Vulnerability", "Function"}, // Example tags
	}

	// Create the sidebar tree widget from the defined sections
	// The tree provides a collapsible, hierarchical view of all items
	tree := createSidebarTree([]SidebarSection{recentNotes, projects, tags})

	// Add padding around the sidebar for visual comfort
	// This creates space between the sidebar elements and the container edges
	return container.NewPadded(tree)
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
