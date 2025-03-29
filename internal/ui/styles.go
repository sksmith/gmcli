package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// UI color palette
const (
	colorPrimary   = "#874BFD"
	colorSecondary = "#7D56F4"
	colorTertiary  = "#5A56F4"
	colorError     = "#FF0000"
	colorSuccess   = "#00FF00"
	colorText      = "#FAFAFA"
	colorMuted     = "#626262"
)

// Exported styles for use in the application
var (
	// AppStyle is the main container style
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorPrimary))

	// TitleStyle for main titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorText)).
			Background(lipgloss.Color(colorSecondary)).
			Padding(0, 1).
			MarginBottom(1)

	// InputStyle for text inputs
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(colorTertiary)).
			Padding(0, 0, 0, 1)
)

// Utility functions for rendering styled text
var (
	// RenderError returns error-styled text
	RenderError = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorError)).
			Render

	// RenderSuccess returns success-styled text
	RenderSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSuccess)).
			Render

	// RenderMuted returns muted-styled text
	RenderMuted = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted)).
			Render
)

// ListStyles returns styles for list components
func ListStyles() (lipgloss.Style, lipgloss.Style, lipgloss.Style) {
	listTitle := TitleStyle.Copy()

	listItem := lipgloss.NewStyle().
		PaddingLeft(2)

	listSelectedItem := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorPrimary)).
		PaddingLeft(2)

	return listTitle, listItem, listSelectedItem
}
