package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sksmith/gmcli/internal/config"
)

// ValidateCalendarName validates a calendar name
func ValidateCalendarName(name string) error {
	if name == "" {
		return fmt.Errorf("calendar name cannot be empty")
	}
	return nil
}

// ValidateCalendarAbbreviation validates a calendar abbreviation
func ValidateCalendarAbbreviation(abbr string) error {
	if abbr == "" {
		return fmt.Errorf("abbreviation cannot be empty")
	}
	if len(abbr) > 3 {
		return fmt.Errorf("abbreviation must be 1-3 characters")
	}
	return nil
}

// ValidateYear validates a year input
func ValidateYear(input string) (int, error) {
	if input == "" {
		return 0, fmt.Errorf("year cannot be empty")
	}

	year, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("year must be a number")
	}

	return year, nil
}

// CreateCalendar creates a new calendar and adds it to the config
func CreateCalendar(cfg *config.Config, input config.CreateCalendarInput) error {
	// Create a new calendar with the provided input
	newCalendar := config.Calendar{
		Name:         input.Name,
		Abbreviation: input.Abbreviation,
		StartYear:    input.StartYear,
		TotalYears:   input.TotalYears,
		Ages:         []config.Age{},
		Months:       []config.Month{},
	}

	// Set default days in year if not set
	if cfg.DaysInYear == 0 {
		cfg.DaysInYear = 365 // Default
	}

	// Create a default age
	defaultAge := config.Age{
		Name:         "First Age",
		Abbreviation: "FA",
		Length:       input.TotalYears,
	}
	newCalendar.Ages = append(newCalendar.Ages, defaultAge)

	// Create 12 months of 30 days each (plus remaining days)
	daysLeft := cfg.DaysInYear
	for i := 1; i <= 12; i++ {
		days := 30
		if i == 12 {
			days = daysLeft // Last month gets remaining days
		}
		daysLeft -= days

		newCalendar.Months = append(newCalendar.Months, config.Month{
			Name: fmt.Sprintf("Month %d", i),
			Days: days,
		})
	}

	// Add the new calendar to config
	cfg.Calendars = append(cfg.Calendars, newCalendar)

	// Save the updated config
	return config.Save(*cfg)
}

// GetCalendarDetails returns a formatted string with calendar details
func GetCalendarDetails(cal config.Calendar) string {
	var details strings.Builder

	details.WriteString(fmt.Sprintf("Calendar: %s (%s)\n", cal.Name, cal.Abbreviation))
	details.WriteString(fmt.Sprintf("Start Year: %d, Total Years: %d\n\n", cal.StartYear, cal.TotalYears))

	details.WriteString("Ages:\n")
	for _, age := range cal.Ages {
		details.WriteString(fmt.Sprintf("- %s (%s): %d years\n",
			age.Name, age.Abbreviation, age.Length))
	}

	details.WriteString("\nMonths:\n")
	for _, month := range cal.Months {
		details.WriteString(fmt.Sprintf("- %s: %d days\n",
			month.Name, month.Days))
	}

	return details.String()
}
