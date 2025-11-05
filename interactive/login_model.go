package interactive

import (
	"fmt"

	"github.com/Eventid3/sql-archiver/mssql"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// -------- FORM MODEL --------------
type loginModel struct {
	focusIndex int
	inputs     []textinput.Model
	err        error
}

func NewLoginModel(err error) loginModel {
	m := loginModel{
		focusIndex: 0,
		inputs:     make([]textinput.Model, 3),
		err:        err,
	}
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "container"
	m.inputs[0].Width = 50
	m.inputs[0].Focus()

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "sa"
	m.inputs[1].Width = 50

	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "password"
	m.inputs[2].Width = 50
	m.inputs[2].EchoMode = textinput.EchoPassword

	return m
}

// Init implements tea.Model.
func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// check for enter on last input
			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				return m, func() tea.Msg {
					err := mssql.CheckConnection(m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value())
					if err != nil {
						return loginFailedMsg{err: err}
					}

					return loginDoneMsg{
						container: m.inputs[0].Value(),
						user:      m.inputs[1].Value(),
						password:  m.inputs[2].Value(),
					}
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			// truncate index if out of bounds
			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					// fm.inputs[i].PromptStyle = focusedStyle
					// fm.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				// fm.inputs[i].PromptStyle = noStyle
				// fm.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *loginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// View implements tea.Model.
func (m loginModel) View() string {
	result := fmt.Sprintf(
		`Fill the form below:

Container: %s
User:      %s
Password:  %s

Enter to submit, Esc or ctrl+q to quit.`,

		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
	)

	if m.err != nil {
		result += errorTextStyle.Render(fmt.Sprintf("\n\nError: %s", m.err.Error()))
	}

	return result
}
