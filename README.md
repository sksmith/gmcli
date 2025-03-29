# Fantasy Calendar CLI

A command-line application for creating and managing fantasy calendars and events.

## Features

- Create and manage fantasy calendars with customizable age systems
- Add events to your calendars
- Track events with their own markdown files
- User-friendly TUI with keyboard navigation

## Installation

```bash
# Clone the repository
git clone https://github.com/sksmith/gmcli.git
cd gmcli

# Build the application
go build -o gmcli ./cmd/gmcli

# Run the application
./gmcli
```

## Usage

The application provides a terminal user interface with the following main features:

- **Create Calendar**: Define a new fantasy calendar with customizable years and date format
- **Create Event**: Add an event to an existing calendar
- **View Calendars**: Browse and inspect your existing calendars
- **Exit**: Close the application

### Navigation

- **↑/k**: Move up
- **↓/j**: Move down
- **Enter**: Select menu item
- **Esc**: Go back
- **?**: Toggle help
- **q**: Quit
- **Ctrl+C**: Force quit

## Directory Structure

- `/templates`: Contains markdown templates for events
- `/events`: Stores generated event files
- `config.yaml`: Application configuration file

## License

This project is licensed under the MIT License - see the LICENSE file for details.