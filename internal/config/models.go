package config

// Config represents the overall configuration.
type Config struct {
	DaysInYear int        `yaml:"days_in_year"`
	Calendars  []Calendar `yaml:"calendars"`
}

// Calendar represents one fantasy calendar.
type Calendar struct {
	Name         string  `yaml:"name"`
	Abbreviation string  `yaml:"abbreviation"`
	StartYear    int     `yaml:"start_year"`
	TotalYears   int     `yaml:"total_years"` // total years available for ages
	Ages         []Age   `yaml:"ages"`
	Months       []Month `yaml:"months"`
}

// Age represents an age in the calendar.
type Age struct {
	Name         string `yaml:"name"`
	Abbreviation string `yaml:"abbreviation"`
	Length       int    `yaml:"length"`                 // in years
	Previous     string `yaml:"previous_age,omitempty"` // optional previous age abbreviation
}

// Month represents one month in the calendar.
type Month struct {
	Name     string `yaml:"name"`
	Days     int    `yaml:"days"`
	Previous string `yaml:"previous_month,omitempty"`
}

// Event represents an event to be created.
type Event struct {
	CalendarName   string
	CalendarAbbrev string
	AgeAbbrev      string
	Year           int
	Month          int
	Day            int
	DaysSinceZero  int
	Name           string
}

// CreateCalendarInput holds data for calendar creation
type CreateCalendarInput struct {
	Name         string
	Abbreviation string
	StartYear    int
	TotalYears   int
}
