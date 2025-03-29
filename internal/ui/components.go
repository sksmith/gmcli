package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/sksmith/gmcli/internal/config"
)

// Item implements list.Item interface for menu options
type Item struct {
	Title       string
	Description string
}

// FilterValue returns the item's title for filtering
func (i Item) FilterValue() string { return i.Title }

// NewTextInput creates a new configured text input
func NewTextInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Focus()
	return input
}

// NewMenuList creates a new list model for menu items
func NewMenuList(items []list.Item, title string, width, height int) list.Model {
	// Configure list delegate
	delegate := list.NewDefaultDelegate()
	listTitle, _, selectedItemStyle := ListStyles()

	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = selectedItemStyle

	// Create the list
	l := list.New(items, delegate, width, height)
	l.Title = title
	l.Styles.Title = listTitle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return l
}

// MainMenuItems returns the items for the main menu
func MainMenuItems() []list.Item {
	return []list.Item{
		Item{Title: "Create Calendar", Description: "Create a new fantasy calendar"},
		Item{Title: "Create Event", Description: "Add an event to an existing calendar"},
		Item{Title: "View Calendars", Description: "View all configured calendars"},
		Item{Title: "Exit", Description: "Exit the application"},
	}
}

// CalendarListItems creates list items from calendars
func CalendarListItems(calendars []config.Calendar) []list.Item {
	items := make([]list.Item, len(calendars))
	for i, cal := range calendars {
		items[i] = Item{
			Title:       cal.Name,
			Description: fmt.Sprintf("Abbreviation: %s, Years: %d", cal.Abbreviation, cal.TotalYears),
		}
	}
	return items
}
