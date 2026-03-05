package main

import (
	"daily_tracker/table"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

type Model struct {
	users         Users
	cursor        int
	modelFilePath string
	table         table.ITable

	input      textinput.Model
	datepicker datepicker.Model

	currentScreen screen
	err           error
}

type Users []User

type User struct {
	Name          string      `json:"name"`
	HostedDailies []time.Time `json:"hosted_dailies"`
}

type screen int

const (
	screenOverview screen = iota
	screenAddUser
	screenUserRename
	screenAddDate
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	helpStyle  = lipgloss.NewStyle().Faint(true)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentScreen {
	case screenOverview:
		return m.updateOverview(msg)
	case screenAddUser:
		return m.updateUserAdd(msg)
	case screenAddDate:
		return m.updateDateAdd(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	switch m.currentScreen {
	case screenOverview:
		return m.viewOverview()
	case screenAddUser:
		return m.viewUserAdd()
	case screenUserRename:
		return "User Rename Screen"
	case screenAddDate:
		return m.viewDateAdd()
	default:
		return "Unknown Screen"
	}
}

// overview screen

func (m Model) updateOverview(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if len(m.users) > 0 {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.users) - 1
				}
			}
		case "down", "j":
			if len(m.users) > 0 {
				if m.cursor < len(m.users)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}
		case "a":
			m.currentScreen = screenAddUser
			m.input.SetValue("")
			m.input.Focus()
			return m, nil
		case "enter":
			m.currentScreen = screenAddDate
			m.datepicker = datepicker.New(time.Now())
			m.datepicker.SelectDate()
			return m, nil
		}
	}

	return m, nil
}

func (m Model) viewOverview() string {
	sb := strings.Builder{}

	if m.err != nil {
		sb.WriteString(errorStyle.Render("⚠ " + m.err.Error()))
		sb.WriteString("\n\n")
	}
	if len(m.users) == 0 {
		sb.WriteString("No users yet. Press 'a' to add one.\n")
	} else {
		sb.WriteString(m.table.ToString())

		str := sb.String()
		lines := strings.Split(str, "\n")
		name := m.users[m.cursor].Name

		for i, line := range lines {
			if strings.Contains(line, name) {
				lines[i] = titleStyle.Render(line)
				break
			}
		}

		sb.Reset()
		sb.WriteString(strings.Join(lines, "\n"))
	}

	// The footer
	sb.WriteString("\nKeys: ↑/↓ move | a add | Enter Add Date | q quit\n")

	// Send the UI for rendering
	return sb.String()
}

// add user screen

func (m Model) updateUserAdd(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.currentScreen = screenOverview
			m.input.Blur()
			return m, tea.ClearScreen

		case "enter":
			name := strings.TrimSpace(m.input.Value())
			if name == "" {
				m.err = errors.New("name cannot be empty")
				return m, nil
			}
			m.users = append(m.users, User{Name: name})
			m.cursor = len(m.users) - 1

			m.err = m.writeChangesToFS()
			m.currentScreen = screenOverview
			m.input.Blur()
			m.updateTable()
			return m, tea.ClearScreen
		}
	}

	return m, cmd
}

func (m Model) viewUserAdd() string {
	var b strings.Builder
	title := "Add user"
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")
	b.WriteString("Name:\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Enter save | Esc cancel"))
	b.WriteString("\n")
	return b.String()
}

// add date screen

func (m Model) updateDateAdd(msg tea.Msg) (Model, tea.Cmd) {
	dp, cmd := m.datepicker.Update(msg)
	m.datepicker = dp

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenOverview
			m.input.Blur()
			return m, nil

		case "enter":
			selectedDate := m.datepicker.Time
			if selectedDate.IsZero() {
				m.err = errors.New("please select a date")
				return m, nil
			}

			m.users[m.cursor].HostedDailies = append(m.users[m.cursor].HostedDailies, selectedDate)

			m.err = m.writeChangesToFS()
			m.currentScreen = screenOverview
			m.datepicker.Blur()
			m.table = m.updateTable()
			return m, tea.ClearScreen
		}
	}

	return m, cmd
}

func (m Model) viewDateAdd() string {
	sb := strings.Builder{}
	sb.WriteString(titleStyle.Render("Add Date for "))
	sb.WriteString(titleStyle.Render(m.users[m.cursor].Name))
	sb.WriteString(titleStyle.Render(":\n"))
	sb.WriteString(m.datepicker.View())
	sb.WriteString("\n")
	sb.WriteString("Enter save | Esc/q cancel")

	return sb.String()
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) writeChangesToFS() error {
	modelFile, err := os.Create(m.modelFilePath)
	if err != nil {
		fmt.Printf("Error opening status file: %v\n", err)
		os.Exit(1)
	}
	defer func(modelFile *os.File) {
		err := modelFile.Close()
		if err != nil {
			m.err = fmt.Errorf("error closing status file: %v", err)
		}
	}(modelFile)

	statusJson, err := json.Marshal(m.users)
	if err != nil {
		return err
	}

	_, err = modelFile.Write(statusJson)
	if err != nil {
		return err
	}

	return nil
}

func initialModel(statusFilePath string) Model {
	var users Users
	var modelFile *os.File
	if _, err := os.Stat(statusFilePath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Status file does not exist at path: %s\n", statusFilePath)
		fmt.Println("Creating a new status file...")

		users = Users{}
	} else {
		modelFile, err = os.Open(statusFilePath)
		if err != nil {
			fmt.Printf("Error opening status file: %v\n", err)
			os.Exit(1)
		}
		defer func(modelFile *os.File) {
			err := modelFile.Close()
			if err != nil {
				fmt.Printf("error closing status file: %v", err)
				os.Exit(1)
			}
		}(modelFile)

		if err = json.NewDecoder(modelFile).Decode(&users); err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("Status file is empty, starting with an empty user list.\n")
				users = Users{}
			} else {
				fmt.Printf("Error decoding status file: %v\n", err)
				os.Exit(1)
			}
		}
	}

	ti := textinput.New()
	ti.Placeholder = "Name"
	ti.CharLimit = 60
	ti.Width = 30

	return Model{
		cursor:        0,
		users:         users,
		modelFilePath: statusFilePath,
		table:         setTable(users),

		input: ti,
	}
}

func setTable(users Users) table.ITable {
	t := table.NewTable()
	err := t.SetPadding(2)
	if err != nil {
		fmt.Printf("Error setting table padding: %v\n", err)
		os.Exit(1)
	}
	t.SetHeaderSeparator("-")
	t.SetColumnSeparator("|")
	t.SetHeader("Name", "Hosted Dailies")
	for _, u := range users {
		// Render the row
		t.AddRow(u.Name, strconv.Itoa(len(u.HostedDailies)))
	}
	return &t
}

func (m Model) updateTable() table.ITable {
	m.table.DeleteRows()
	for _, u := range m.users {
		// Render the row
		m.table.AddRow(u.Name, strconv.Itoa(len(u.HostedDailies)))
	}

	return m.table
}

func main() {
	statusFilePath := flag.String("f", "./daily_tracker.json", "Path to the status file")
	printHelp := flag.Bool("h", false, "Print help message")
	flag.Parse()

	if *printHelp {
		flag.Usage()
		os.Exit(0)
	}

	p := tea.NewProgram(initialModel(*statusFilePath), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
