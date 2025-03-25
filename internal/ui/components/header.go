// Package components provides UI components for the RevEnGo application.
// This file contains the header component used at the top of the main window.
package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/leog/RevEnGo/internal/ui/widgets"
)

// Color constants for the header
var (
	headerStartColor = color.NRGBA{R: 15, G: 23, B: 42, A: 255}   // Darker blue
	headerEndColor   = color.NRGBA{R: 30, G: 41, B: 59, A: 255}   // Slightly lighter blue
	accentColor      = color.NRGBA{R: 0, G: 174, B: 239, A: 255}  // Cyber blue for accents
	glowColor        = color.NRGBA{R: 61, G: 134, B: 247, A: 100} // Glowing effect color
)

// NewHeader creates a new header component for the application.
// The header includes:
// - Application title/logo
// - Main navigation buttons
// - User actions (settings, help, etc.)
//
// Returns a canvas object that can be placed in a container.
func NewHeader() fyne.CanvasObject {
	// Create a gradient background for the header
	gradient := widgets.GradientBackground(headerStartColor, headerEndColor)

	// Create the application title with digital styling
	appTitle := widgets.DigitalText("RevEnGo", 22, accentColor)

	// Create action buttons with cyber styling
	newButton := widgets.CyberButton("New", theme.DocumentCreateIcon(), nil)
	openButton := widgets.CyberButton("Open", theme.FolderOpenIcon(), nil)
	saveButton := widgets.CyberButton("Save", theme.DocumentSaveIcon(), nil)

	// Add subtle glow effects to the buttons
	newButtonWithGlow := widgets.GlowEffect(newButton, glowColor)
	openButtonWithGlow := widgets.GlowEffect(openButton, glowColor)
	saveButtonWithGlow := widgets.GlowEffect(saveButton, glowColor)

	// Create status indicator for system health
	statusIndicator := widgets.StatusIndicator("secure")
	statusLabel := canvas.NewText("System Status: Secure", color.NRGBA{R: 60, G: 220, B: 100, A: 255})
	statusLabel.TextSize = 12
	statusContainer := container.NewHBox(statusIndicator, statusLabel)

	// Create a toolbar with the action buttons
	toolbar := container.NewHBox(
		newButtonWithGlow,
		openButtonWithGlow,
		saveButtonWithGlow,
		layout.NewSpacer(),
		statusContainer,
	)

	// Add cybersecurity-themed decorative elements
	// Binary pattern as decoration
	binaryPattern := canvas.NewText("01001010 10101", color.NRGBA{R: 100, G: 150, B: 200, A: 100})
	binaryPattern.TextSize = 10
	binaryPattern.TextStyle = fyne.TextStyle{Monospace: true}

	// Create hexagonal icon indicator
	hexIcon := widgets.HexagonalButton(theme.InfoIcon(), nil)

	// Arrange the title and decorations
	titleContainer := container.NewHBox(
		hexIcon,
		layout.NewSpacer(),
		appTitle,
		layout.NewSpacer(),
		binaryPattern,
	)

	// Arrange the entire header with the gradient background
	content := container.NewVBox(
		titleContainer,
		widget.NewSeparator(),
		toolbar,
	)

	// Combine the gradient and content
	return container.NewStack(
		gradient,
		container.NewPadded(content),
	)
}
