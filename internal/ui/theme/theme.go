// Package theme provides custom theming for the RevEnGo application.
package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Custom colors for the RevEnGo theme
var (
	// Main background color - deep charcoal with slight blue tint
	colorBackground = color.NRGBA{R: 16, G: 20, B: 30, A: 255}

	// Primary action color - cyber blue
	colorPrimary = color.NRGBA{R: 0, G: 174, B: 239, A: 255}

	// Secondary action color - cyan highlight
	colorSecondary = color.NRGBA{R: 80, G: 200, B: 220, A: 255}

	// Accent color - electric blue
	colorAccent = color.NRGBA{R: 61, G: 134, B: 247, A: 255}

	// Terminal green for code elements
	colorTerminal = color.NRGBA{R: 35, G: 209, B: 139, A: 255}

	// Warning color - amber highlight
	colorWarning = color.NRGBA{R: 255, G: 170, B: 0, A: 255}

	// Error color - digital red
	colorError = color.NRGBA{R: 250, G: 70, B: 76, A: 255}

	// Button background - dark with slight glow
	colorButton = color.NRGBA{R: 30, G: 40, B: 60, A: 255}

	// Input background - darker than main background
	colorInputBackground = color.NRGBA{R: 12, G: 15, B: 22, A: 255}
)

// RevEnGoTheme is a custom theme for the RevEnGo application
type RevEnGoTheme struct{}

// New creates a new instance of the RevEnGoTheme
func New() fyne.Theme {
	return &RevEnGoTheme{}
}

// Color returns the color for the specified ColorName and theme
func (t *RevEnGoTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colorBackground
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNamePrimary:
		return colorPrimary
	case theme.ColorNameButton:
		return colorButton
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 40, G: 50, B: 70, A: 200}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 30, G: 40, B: 50, A: 120}
	case theme.ColorNameInputBackground:
		return colorInputBackground
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 100, G: 120, B: 140, A: 200}
	case theme.ColorNameHover:
		return color.NRGBA{R: 60, G: 80, B: 120, A: 30}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 10, G: 120, B: 200, A: 60}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 30, G: 150, B: 220, A: 60}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font returns the font resource for the specified TextStyle and theme
func (t *RevEnGoTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		// Use a more distinct monospace font for code
		return theme.DefaultTheme().Font(style)
	}
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the specified IconName and theme
func (t *RevEnGoTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// For simplicity, we'll use the default icons for now
	// Later we can customize specific icons as needed
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the specified SizeName and theme
func (t *RevEnGoTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 4
	case theme.SizeNameText:
		return 13
	case theme.SizeNameHeadingText:
		return 18
	case theme.SizeNameSubHeadingText:
		return 15
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInlineIcon:
		return 20
	default:
		return theme.DefaultTheme().Size(name)
	}
}
