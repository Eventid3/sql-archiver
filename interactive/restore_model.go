package interactive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type restoreModel struct {
	focusIndex                 int
	inputs                     []textinput.Model
	filename, mdfFile, ldfFile string
	// config           ServerConfig
}

func NewRestoreModel(config ServerConfig, filename, mdf, ldf string) restoreModel {
	m := restoreModel{
		focusIndex: 0,
		inputs:     make([]textinput.Model, 1),
		filename:   filename,
		mdfFile:    mdf,
		ldfFile:    ldf,
		// config:     config,
	}

	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = mdf
	m.inputs[0].Width = 50
	m.inputs[0].Focus()

	return m
}

func (m restoreModel) Init() tea.Cmd {
	return nil
}

func (m restoreModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					return restoreExecMsg{
						m.filename,
						m.mdfFile,
						m.ldfFile,
						m.inputs[0].Value(),
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

func (m *restoreModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m restoreModel) View() string {
	return fmt.Sprintf(
		`
New database name: %s

Press enter to confirm restore, Esc or ctrl+q to cancel.
		`,
		m.inputs[0].View(),
	)
}
