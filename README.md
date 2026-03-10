# Daily Tracker

A simple TUI (Terminal User Interface) application for tracking daily standups and who hosted them. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **User Management**: Add and manage team members
- **Daily Tracking**: Record when someone hosted a daily
- **Clear Overview**: Table view of all users with the number of hosted dailies
- **Persistence**: Data is saved to a JSON file

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd daily_tracker
```
```bash
# Install dependencies
go mod download
```
```bash
# Run the application
go run .
```
```bas
# build executable
go build .
```

## Usage

```bash
# Start with default file (./daily_tracker.json)
./daily_tracker

# Start with custom status file
./daily_tracker -f /path/to/file.json

# Show help
./daily_tracker -h
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-f` | Path to status file | `./daily_tracker.json` |
| `-h` | Show help | - |

## Keyboard Shortcuts

### Overview Screen

| Key | Action |
|-----|--------|
| `↑` / `k` | Navigate up |
| `↓` / `j` | Navigate down |
| `a` | Add new user |
| `Enter` | Add date for selected user |
| `q` / `Ctrl+C` | Quit |

### Add User

| Key | Action |
|-----|--------|
| `Enter` | Save |
| `Esc` | Cancel |

### Add Date

| Key | Action |
|-----|--------|
| `Enter` | Save date |
| `Esc` / `q` | Cancel |

## Data Format

The application stores data in JSON format:

```json
[
  {
    "name": "Max Mustermann",
    "hosted_dailies": [
      "2026-03-01T10:00:00Z",
      "2026-03-05T10:00:00Z"
    ]
  }
]
```

## Development

### Run Tests

```bash
go test ./...
```

### Project Structure

```
daily_tracker/
├── main.go           # Main application with TUI logic
├── main_test.go      # Tests for the main application
├── go.mod            # Go Module definition
├── go.sum            # Dependency checksums
├── daily_tracker.json # Data file (created automatically)
└── table/
    ├── table.go      # Table rendering logic
    └── table_test.go # Tests for the table
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI Framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI Components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style Definitions
- [bubble-datepicker](https://github.com/ethanefung/bubble-datepicker) - Date Picker

## Licence
GNU GPL v3
