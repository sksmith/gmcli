package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/sksmith/gmcli/internal/config"
)

// ValidateEventDate validates an event date in format AAYYYY-MM-DD
func ValidateEventDate(dateStr string, cal config.Calendar, daysInYear int) (config.Event, error) {
	// Create empty event for returning errors
	emptyEvent := config.Event{}

	// Check basic format
	if len(dateStr) != 12 || dateStr[6] != '-' || dateStr[9] != '-' {
		return emptyEvent, fmt.Errorf("invalid format, date must be AAYYYY-MM-DD")
	}

	// Extract components
	ageAbbrev := dateStr[:2]
	yearStr := dateStr[2:6]
	monthStr := dateStr[7:9]
	dayStr := dateStr[10:12]

	// Validate components can be parsed as numbers
	year, err1 := strconv.Atoi(yearStr)
	month, err2 := strconv.Atoi(monthStr)
	day, err3 := strconv.Atoi(dayStr)

	if err1 != nil || err2 != nil || err3 != nil {
		return emptyEvent, fmt.Errorf("date contains invalid numbers")
	}

	// Check if age abbreviation exists
	foundAge := false
	for _, a := range cal.Ages {
		if a.Abbreviation == ageAbbrev {
			foundAge = true
			break
		}
	}
	if !foundAge {
		return emptyEvent, fmt.Errorf("age abbreviation '%s' not found in calendar", ageAbbrev)
	}

	// Validate month
	if month < 1 || month > len(cal.Months) {
		return emptyEvent, fmt.Errorf("month must be between 1 and %d", len(cal.Months))
	}

	// Validate day
	selectedMonth := cal.Months[month-1]
	if day < 1 || day > selectedMonth.Days {
		return emptyEvent, fmt.Errorf("day must be between 1 and %d for month '%s'",
			selectedMonth.Days, selectedMonth.Name)
	}

	// Calculate days since 0
	totalDays := year * daysInYear

	// Sum days of months before the selected month
	for i := 0; i < month-1; i++ {
		totalDays += cal.Months[i].Days
	}
	totalDays += day

	// Create and return the event
	event := config.Event{
		CalendarName:   cal.Name,
		CalendarAbbrev: cal.Abbreviation,
		AgeAbbrev:      ageAbbrev,
		Year:           year,
		Month:          month,
		Day:            day,
		DaysSinceZero:  totalDays,
	}

	return event, nil
}

// ValidateEventName validates an event name
func ValidateEventName(name string) error {
	if name == "" {
		return fmt.Errorf("event name cannot be empty")
	}
	return nil
}

// CreateEvent creates a new event from the provided data
func CreateEvent(cal config.Calendar, event config.Event) error {
	// Load template
	tmplPath := filepath.Join("templates", "event.md.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Create a file name based on event name and daysSinceZero
	safeName := strings.ReplaceAll(strings.ToLower(event.Name), " ", "_")
	outFile := filepath.Join("events", fmt.Sprintf("%s_%d.md", safeName, event.DaysSinceZero))

	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create event file: %w", err)
	}
	defer f.Close()

	// Execute template
	if err := tmpl.Execute(f, event); err != nil {
		return fmt.Errorf("failed to write event file: %w", err)
	}

	return nil
}
