package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderCalendarPreview returns a formatted string with calendar info
func RenderCalendarPreview(name, abbrev string, startYear, totalYears int) string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("Calendar Preview"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Name: %s\n", name))
	b.WriteString(fmt.Sprintf("Abbreviation: %s\n", abbrev))
	b.WriteString(fmt.Sprintf("Start Year: %d\n", startYear))
	b.WriteString(fmt.Sprintf("Total Years: %d\n", totalYears))

	return b.String()
}

// RenderEventPreview returns a formatted string with event info
func RenderEventPreview(calendarName string, date string, daysSinceZero int) string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("Event Preview"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Calendar: %s\n", calendarName))
	b.WriteString(fmt.Sprintf("Date: %s\n", date))
	b.WriteString(fmt.Sprintf("Days Since Year 0: %d\n", daysSinceZero))

	return b.String()
}

// RenderConfigPreview returns a formatted string with configuration preview
func RenderConfigPreview(daysInYear int, calendarCount int) string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("Configuration Preview"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Days in Year: %d\n", daysInYear))
	b.WriteString(fmt.Sprintf("Calendars: %d\n", calendarCount))

	return b.String()
}

// JoinHorizontal joins content with a divider
func JoinHorizontal(left, right string) string {
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		SetString("│").
		String()

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		divider,
		right,
	)
}
