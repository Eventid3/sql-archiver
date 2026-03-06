package interactive

import (
	"fmt"
	"strings"

	"github.com/Eventid3/sql-archiver/config"
	"github.com/Eventid3/sql-archiver/domain"
	"github.com/Eventid3/sql-archiver/mssql"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// loginFocus tracks which section of the login screen has keyboard focus.
type loginFocus int

const (
	loginFocusPicker loginFocus = iota // saved-connection list
	loginFocusForm                     // host/user/password inputs
)

// loginModel shows a saved-connection picker above the manual login form.
//
// Keyboard map (picker focused):
//
//	↑/↓        navigate saved connections
//	Enter      pre-fill the form with the selected connection and shift focus to the form
//	n          clear the form and shift focus to the form (new connection)
//	d          show inline delete confirmation for the highlighted connection
//	y/N        confirm or cancel the pending delete
//
// Keyboard map (form focused):
//
//	Tab/↓/Enter   advance focus to the next field
//	Shift+Tab/↑   move focus to the previous field
//	Enter (last)  validate + attempt connection
//	Esc           return focus to the picker (or quit if no saved connections)
type loginModel struct {
	// saved connections shown in the picker
	connections []domain.Connection
	pickerIndex int // currently highlighted row in the picker

	// form fields
	focusIndex int
	inputs     []textinput.Model

	// which section is active
	focus loginFocus

	// inline delete confirmation state
	pendingDelete *domain.Connection // non-nil while awaiting y/N

	err error
}

// NewLoginModel constructs a login screen.  connections is the list loaded from
// the config file; err (if non-nil) is displayed as an error banner.
func NewLoginModel(connections []domain.Connection, err error) loginModel {
	m := loginModel{
		connections: connections,
		focusIndex:  0,
		inputs:      make([]textinput.Model, 3),
		err:         err,
	}

	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "localhost:1433"
	m.inputs[0].Width = 50

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "sa"
	m.inputs[1].Width = 50

	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "password"
	m.inputs[2].Width = 50
	m.inputs[2].EchoMode = textinput.EchoPassword

	if len(connections) > 0 {
		m.focus = loginFocusPicker
	} else {
		m.focus = loginFocusForm
		m.inputs[0].Focus()
	}
	return m
}

// prefillForm loads the credentials of a saved connection into the form inputs.
func (m *loginModel) prefillForm(c domain.Connection) {
	m.inputs[0].SetValue(c.Host)
	m.inputs[1].SetValue(c.User)
	m.inputs[2].SetValue(c.Password)
}

// clearForm resets all form inputs.
func (m *loginModel) clearForm() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
}

// shiftFocusToForm moves keyboard focus to the form and focuses the first field.
func (m *loginModel) shiftFocusToForm() tea.Cmd {
	m.focus = loginFocusForm
	m.focusIndex = 0
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	return m.inputs[0].Focus()
}

// Init implements tea.Model.
func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Inline delete confirmation takes priority.
		if m.pendingDelete != nil {
			return m.handleDeleteConfirm(msg)
		}

		if m.focus == loginFocusPicker {
			return m.handlePickerKey(msg)
		}
		return m.handleFormKey(msg)
	}

	// Forward character input to whichever text inputs are focused.
	cmd := m.updateInputs(msg)
	return m, cmd
}

// handlePickerKey handles key events while the saved-connection picker is focused.
func (m loginModel) handlePickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "shift+tab":
		if m.pickerIndex > 0 {
			m.pickerIndex--
		}
		return m, nil

	case "down", "tab":
		// +1 for the "New connection" row
		if m.pickerIndex < len(m.connections) {
			m.pickerIndex++
		}
		return m, nil

	case "enter":
		if m.pickerIndex == len(m.connections) {
			// "New connection" row selected
			m.clearForm()
			cmd := m.shiftFocusToForm()
			return m, cmd
		}
		m.prefillForm(m.connections[m.pickerIndex])
		host := m.inputs[0].Value()
		user := m.inputs[1].Value()
		password := m.inputs[2].Value()
		return m, func() tea.Msg {
			if err := mssql.CheckConnection(host, user, password); err != nil {
				return loginFailedMsg{err: err}
			}
			return loginDoneMsg{host: host, user: user, password: password}
		}

	case "n":
		m.clearForm()
		cmd := m.shiftFocusToForm()
		return m, cmd

	case "d":
		if m.pickerIndex < len(m.connections) {
			c := m.connections[m.pickerIndex]
			m.pendingDelete = &c
		}
		return m, nil

	case "esc":
		return m, tea.Quit
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

// handleDeleteConfirm handles the y/N confirmation for deleting a connection.
func (m loginModel) handleDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "y":
		name := m.pendingDelete.Name
		m.pendingDelete = nil
		return m, func() tea.Msg {
			if err := config.DeleteConnection(name); err != nil {
				return loginFailedMsg{err: err}
			}
			return connectionDeletedMsg{}
		}
	default:
		// Any other key cancels.
		m.pendingDelete = nil
		return m, nil
	}
}

// handleFormKey handles key events while the manual-login form is focused.
func (m loginModel) handleFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		if len(m.connections) > 0 {
			// Return to picker
			m.focus = loginFocusPicker
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, nil
		}
		return m, tea.Quit

	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		// Submit on Enter at the last field.
		if s == "enter" && m.focusIndex == len(m.inputs)-1 {
			host := m.inputs[0].Value()
			user := m.inputs[1].Value()
			password := m.inputs[2].Value()
			savedConns := m.connections
			return m, func() tea.Msg {
				if err := mssql.CheckConnection(host, user, password); err != nil {
					return loginFailedMsg{err: err}
				}
				if host == "" {
					host = "localhost:1433"
				}
				// Check whether this host+user is already saved.
				for _, c := range savedConns {
					if c.Host == host && c.User == user {
						return loginDoneMsg{host: host, user: user, password: password}
					}
				}
				return promptSaveConnectionMsg{host: host, user: user, password: password}
			}
		}

		// Move focus between fields.
		if s == "up" || s == "shift+tab" {
			m.focusIndex--
		} else {
			m.focusIndex++
		}
		if m.focusIndex > len(m.inputs)-1 {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		cmds[m.focusIndex] = m.inputs[m.focusIndex].Focus()
		return m, tea.Batch(cmds...)
	}

	// Forward all other key events (character input) to the active text input.
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *loginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m loginModel) View() string {
	var b strings.Builder

	// ── Saved connections picker ──────────────────────────────────────────
	if len(m.connections) > 0 {
		b.WriteString("Saved connections  (↑↓ navigate · Enter select · n new · d delete)\n\n")
		for i, c := range m.connections {
			cursor := "  "
			label := fmt.Sprintf("%-20s  %s@%s", c.Name, c.User, c.Host)
			if m.focus == loginFocusPicker && i == m.pickerIndex {
				cursor = "> "
				label = ColHeaderStyle.Render(label)
			}
			if m.pendingDelete != nil && c.Name == m.pendingDelete.Name {
				label += ErrorTextStyle.Render(fmt.Sprintf("  Delete '%s'? [y/N]", c.Name))
			}
			b.WriteString(cursor + label + "\n")
		}
		// "New connection" row
		newLabel := fmt.Sprintf("  %-20s", "+ New connection")
		if m.focus == loginFocusPicker && m.pickerIndex == len(m.connections) {
			newLabel = "> " + ColHeaderStyle.Render("+ New connection")
		}
		b.WriteString(newLabel + "\n")
		b.WriteString("\n")
	}

	// ── Manual login form ─────────────────────────────────────────────────
	formSection := fmt.Sprintf(
		`Fill the form below:

Host:      %s
User:      %s
Password:  %s

Enter to submit · Esc to go back`,
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
	)

	if m.focus == loginFocusForm {
		b.WriteString(formSection)
	} else {
		b.WriteString(DocStyle.Faint(true).Render(formSection))
	}

	if m.err != nil {
		b.WriteString(ErrorTextStyle.Render(fmt.Sprintf("\n\nError: %s", m.err.Error())))
	}

	return b.String()
}
