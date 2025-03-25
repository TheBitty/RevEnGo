// Package widgets provides custom UI widgets for the RevEnGo application.
package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CyberButton creates a custom styled button with a cybersecurity aesthetic
func CyberButton(text string, icon fyne.Resource, onTapped func()) fyne.CanvasObject {
	// Create the base button
	btn := widget.NewButtonWithIcon(text, icon, onTapped)

	// Apply custom styling to make it look more cyber-themed
	btn.Importance = widget.MediumImportance

	return btn
}

// TerminalEntry creates a text entry field styled like a terminal
func TerminalEntry() *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.TextStyle = fyne.TextStyle{Monospace: true}

	return entry
}

// HexTabItem creates a tab item styled like a memory address
func HexTabItem(title string, content fyne.CanvasObject) *container.TabItem {
	// Format the title to look like a hex address
	hexTitle := "0x" + title

	return container.NewTabItem(hexTitle, content)
}

// GradientBackground creates a gradient background with specified colors
func GradientBackground(startColor, endColor color.Color) *canvas.LinearGradient {
	gradient := canvas.NewLinearGradient(startColor, endColor, 0)
	return gradient
}

// GlowEffect adds a subtle glow effect to a canvas object
func GlowEffect(obj fyne.CanvasObject, glowColor color.Color) fyne.CanvasObject {
	// Create a rectangle with the glow color
	glowRect := canvas.NewRectangle(glowColor)

	// Overlay the original object on top of the glow
	return container.NewStack(glowRect, obj)
}

// HexagonalButton creates a hexagonal-shaped button
// Note: Fyne doesn't directly support custom shapes, so this is approximate
func HexagonalButton(icon fyne.Resource, onTapped func()) fyne.CanvasObject {
	btn := widget.NewButtonWithIcon("", icon, onTapped)

	// Apply styling to make it look more distinct
	btn.Importance = widget.HighImportance

	return btn
}

// DigitalText creates text with a digital/circuit styling
func DigitalText(text string, size float32, textColor color.Color) fyne.CanvasObject {
	label := canvas.NewText(text, textColor)
	label.TextSize = size
	label.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	return label
}

// StatusIndicator creates a colored circle indicator for status
func StatusIndicator(status string) fyne.CanvasObject {
	var statusColor color.Color

	// Choose color based on status
	switch status {
	case "vulnerable":
		statusColor = color.NRGBA{R: 255, G: 60, B: 60, A: 255} // Red
	case "secure":
		statusColor = color.NRGBA{R: 60, G: 220, B: 100, A: 255} // Green
	case "warning":
		statusColor = color.NRGBA{R: 255, G: 180, B: 50, A: 255} // Yellow
	case "info":
		statusColor = color.NRGBA{R: 80, G: 170, B: 255, A: 255} // Blue
	default:
		statusColor = color.Gray{Y: 150} // Gray for unknown
	}

	// Create a circle with the status color
	circle := canvas.NewCircle(statusColor)
	circle.StrokeWidth = 1
	circle.StrokeColor = color.White

	// Set a fixed size for the circle
	circleContainer := container.NewWithoutLayout(circle)
	circleContainer.Resize(fyne.NewSize(12, 12))

	return circleContainer
}

// NotePadWithLineNumbers creates a notepad with line numbers for code
func NotePadWithLineNumbers() fyne.CanvasObject {
	// Create the main text area
	codeArea := widget.NewMultiLineEntry()
	codeArea.TextStyle = fyne.TextStyle{Monospace: true}

	// Create line numbers (simplified version, not dynamic)
	lineNumbers := widget.NewLabel("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")
	lineNumbers.TextStyle = fyne.TextStyle{Monospace: true}
	lineNumbers.Alignment = fyne.TextAlignTrailing

	// Combine them in a horizontal container
	return container.NewBorder(nil, nil, lineNumbers, nil, codeArea)
}

// HexDumpView creates a view that mimics a hex editor
func HexDumpView() fyne.CanvasObject {
	// Create the hex dump text
	hexText := widget.NewLabel("00000000  2e 74 65 78 74 00 00 00  00 00 00 00 00 00 00 00  |.text...........|")
	hexText.TextStyle = fyne.TextStyle{Monospace: true}

	// Add more rows for demonstration
	hexRows := []string{
		"00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|",
		"00000020  2e 64 61 74 61 00 00 00  00 00 00 00 00 00 00 00  |.data...........|",
		"00000030  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|",
	}

	for _, row := range hexRows {
		hexText.SetText(hexText.Text + "\n" + row)
	}

	// Create a scroll container for the hex dump
	return container.NewScroll(hexText)
}
