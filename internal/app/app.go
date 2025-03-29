package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sksmith/gmcli/internal/commands"
	"github.com/sksmith/gmcli/internal/config"
	"github.com/sksmith/gmcli/internal/ui"
)

// State constants
const (
	stateMenu           = "menu"
	stateCreateCalendar = "create_calendar"
	stateSelectCalendar = "select_calendar"
	stateViewCalendars  = "view_calendars"
	stateEventDate      = "event_date"
	stateEventName      = "event_name"
)

// AppModel represents the application state
type AppModel struct {
	config        config.Config
	state         string
	width, height int
	input         textinput.Model
	menuList      list.Model
	keymap        KeyMap
	helper        help.Model
	statusMsg     string
	helpEnabled   bool

	// Create Calendar fields
	calendarInput      config.CreateCalendarInput
	calendarInputStage int

	// Create Event fields
	eventCalendarIndex int
	eventData          config.Event
	eventDateStr       string
}

// Start initializes and runs the application
func Start() error {
	// Initialize app and configuration
	if err := config.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	program := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := program.Run()
	return err
}

// initialModel creates the initial application model
func initialModel() AppModel {
	var statusMsg string

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		statusMsg = ui.RenderError(fmt.Sprintf("Error loading config: %v", err))
	} else if len(cfg.Calendars) > 0 {
		statusMsg = ui.RenderSuccess("Configuration loaded successfully.")
	} else {
		statusMsg = "No existing configuration found. Starting fresh."
	}

	// Create menu list
	menuList := ui.NewMenuList(ui.MainMenuItems(), "Fantasy Calendar CLI", 0, 0)

	// Return the model
	return AppModel{
		state:       stateMenu,
		config:      cfg,
		input:       ui.NewTextInput(""),
		menuList:    menuList,
		keymap:      DefaultKeyMap(),
		helper:      help.New(),
		statusMsg:   statusMsg,
		helpEnabled: false,
	}
}

// Init initializes the model
func (m AppModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// View renders the current UI state
func (m AppModel) View() string {
	var content string

	switch m.state {
	case stateMenu, stateSelectCalendar, stateViewCalendars:
		content = m.menuList.View()
	case stateCreateCalendar, stateEventDate, stateEventName:
		var headerText string
		switch m.state {
		case stateCreateCalendar:
			switch m.calendarInputStage {
			case 1:
				headerText = "Create Calendar - Name"
			case 2:
				headerText = "Create Calendar - Abbreviation"
			case 3:
				headerText = "Create Calendar - Start Year"
			case 4:
				headerText = "Create Calendar - Total Years"
			}
		case stateEventDate:
			headerText = "Create Event - Enter Date"
		case stateEventName:
			headerText = "Create Event - Enter Name"
		}

		header := ui.TitleStyle.Render(headerText)
		content = lipgloss.JoinVertical(lipgloss.Left,
			header,
			"",
			m.input.View(),
		)
	}

	return ui.AppStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.header(),
			content,
			m.footer(),
		),
	)
}

// header returns the app header
func (m AppModel) header() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		ui.TitleStyle.Render("Fantasy Calendar CLI"),
		"",
	)
}

// footer returns the app footer with help and status messages
func (m AppModel) footer() string {
	var footer strings.Builder

	// Add help text if enabled
	if m.helpEnabled {
		footer.WriteString("\n" + m.helper.View(m.keymap))
	}

	// Add status message if present
	if m.statusMsg != "" {
		if footer.Len() > 0 {
			footer.WriteString("\n\n")
		} else {
			footer.WriteString("\n")
		}
		footer.WriteString(m.statusMsg)
	}

	// Add basic help hint if no help is shown
	if !m.helpEnabled && footer.Len() == 0 {
		footer.WriteString("\n" + ui.RenderMuted("Press ? for help"))
	}

	return footer.String()
}

// Update handles messages and state transitions
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		headerHeight := lipgloss.Height(m.header())
		footerHeight := lipgloss.Height(m.footer())
		verticalMarginHeight := headerHeight + footerHeight

		// Update list height
		h, v := ui.AppStyle.GetFrameSize()
		m.menuList.SetSize(msg.Width-h, msg.Height-v-verticalMarginHeight)
		return m, nil

	case tea.KeyMsg:
		// Handle global keybindings first
		switch {
		case key.Matches(msg, m.keymap.ForceQuit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Help):
			m.helpEnabled = !m.helpEnabled
			return m, nil

		case key.Matches(msg, m.keymap.Quit) && m.state == stateMenu:
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Back):
			if m.state != stateMenu {
				m.state = stateMenu
				m.statusMsg = ""
				m.menuList.SetItems(ui.MainMenuItems())
				m.menuList.Title = "Fantasy Calendar CLI"
				return m, nil
			}
		}

		// State-specific handling
		switch m.state {
		case stateMenu:
			// Handle list navigation and selection
			m.menuList, cmd = m.menuList.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, m.keymap.Enter) {
				item, ok := m.menuList.SelectedItem().(ui.Item)
				if ok {
					switch item.Title {
					case "Create Calendar":
						m.state = stateCreateCalendar
						m.input = ui.NewTextInput("Enter calendar name")
						m.calendarInputStage = 1
						m.statusMsg = ""

					case "Create Event":
						if len(m.config.Calendars) == 0 {
							m.statusMsg = ui.RenderError("No calendars available. Create a calendar first.")
						} else {
							m.state = stateSelectCalendar
							m.menuList.SetItems(ui.CalendarListItems(m.config.Calendars))
							m.menuList.Title = "Select Calendar"
							m.statusMsg = ""
						}

					case "View Calendars":
						if len(m.config.Calendars) == 0 {
							m.statusMsg = ui.RenderError("No calendars to view.")
						} else {
							m.state = stateViewCalendars
							m.menuList.SetItems(ui.CalendarListItems(m.config.Calendars))
							m.menuList.Title = "Calendars"
							m.statusMsg = ""
						}

					case "Exit":
						return m, tea.Quit
					}
				}
			}

		case stateViewCalendars:
			// Handle viewing calendar details
			m.menuList, cmd = m.menuList.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, m.keymap.Enter) {
				item, ok := m.menuList.SelectedItem().(ui.Item)
				if ok {
					for _, cal := range m.config.Calendars {
						if cal.Name == item.Title {
							m.statusMsg = commands.GetCalendarDetails(cal)
							break
						}
					}
				}
			}

		case stateSelectCalendar:
			// Handle calendar selection for event creation
			m.menuList, cmd = m.menuList.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, m.keymap.Enter) {
				item, ok := m.menuList.SelectedItem().(ui.Item)
				if ok {
					for idx, cal := range m.config.Calendars {
						if cal.Name == item.Title {
							m.eventCalendarIndex = idx
							m.state = stateEventDate
							m.input = ui.NewTextInput("Format: AAYYYY-MM-DD (e.g., AB0001-01-01)")
							m.statusMsg = ""
							break
						}
					}
				}
			}

		case stateCreateCalendar:
			// Handle calendar creation flow (multi-step)
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, m.keymap.Enter) {
				input := strings.TrimSpace(m.input.Value())

				switch m.calendarInputStage {
				case 1: // Calendar name
					if err := commands.ValidateCalendarName(input); err != nil {
						m.statusMsg = ui.RenderError(err.Error())
						return m, nil
					}
					m.calendarInput.Name = input
					m.input = ui.NewTextInput("Enter abbreviation (1-3 chars)")
					m.calendarInputStage = 2

				case 2: // Calendar abbreviation
					if err := commands.ValidateCalendarAbbreviation(input); err != nil {
						m.statusMsg = ui.RenderError(err.Error())
						return m, nil
					}
					m.calendarInput.Abbreviation = input
					m.input = ui.NewTextInput("Enter start year (number)")
					m.calendarInputStage = 3

				case 3: // Start year
					startYear, err := commands.ValidateYear(input)
					if err != nil {
						m.statusMsg = ui.RenderError(err.Error())
						return m, nil
					}
					m.calendarInput.StartYear = startYear
					m.input = ui.NewTextInput("Enter total years (number)")
					m.calendarInputStage = 4

				case 4: // Total years
					totalYears, err := commands.ValidateYear(input)
					if err != nil {
						m.statusMsg = ui.RenderError(err.Error())
						return m, nil
					}
					m.calendarInput.TotalYears = totalYears

					// Create the calendar
					if err := commands.CreateCalendar(&m.config, m.calendarInput); err != nil {
						m.statusMsg = ui.RenderError(fmt.Sprintf("Failed to create calendar: %v", err))
					} else {
						m.statusMsg = ui.RenderSuccess("Calendar created successfully!")
					}

					// Reset and return to main menu
					m.state = stateMenu
					m.calendarInputStage = 0
					m.menuList.SetItems(ui.MainMenuItems())
					m.menuList.Title = "Fantasy Calendar CLI"
				}
			}

		case stateEventDate:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, m.keymap.Enter) {
				m.eventDateStr = strings.TrimSpace(m.input.Value())

				// Validate the date format
				eventData, err := commands.ValidateEventDate(
					m.eventDateStr,
					m.config.Calendars[m.eventCalendarIndex],
					m.config.DaysInYear)

				if err != nil {
					m.statusMsg = ui.RenderError(err.Error())
					return m, nil
				}

				// Store the event data and move to name input
				m.eventData = eventData
				m.state = stateEventName
				m.input = ui.NewTextInput("Enter event name")
				m.statusMsg = fmt.Sprintf("Event date: %s (Days since 0: %d)",
					m.eventDateStr, eventData.DaysSinceZero)
			}

		case stateEventName:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, m.keymap.Enter) {
				eventName := strings.TrimSpace(m.input.Value())

				if err := commands.ValidateEventName(eventName); err != nil {
					m.statusMsg = ui.RenderError(err.Error())
					return m, nil
				}

				// Update event name
				m.eventData.Name = eventName

				// Create the event
				cal := m.config.Calendars[m.eventCalendarIndex]
				if err := commands.CreateEvent(cal, m.eventData); err != nil {
					m.statusMsg = ui.RenderError(fmt.Sprintf("Failed to create event: %v", err))
				} else {
					m.statusMsg = ui.RenderSuccess(fmt.Sprintf("Event '%s' created successfully!", eventName))
				}

				// Reset and return to main menu
				m.state = stateMenu
				m.menuList.SetItems(ui.MainMenuItems())
				m.menuList.Title = "Fantasy Calendar CLI"
			}
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
