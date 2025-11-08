package interactive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type restoreModel struct {
	focusIndex int
	inputs     []textinput.Model
	// bakFileName string
	bakFileInfo BakFileInfo
}

func NewRestoreModel(config ServerConfig, fileInfo BakFileInfo) restoreModel {
	m := restoreModel{
		focusIndex:  0,
		inputs:      make([]textinput.Model, 1),
		bakFileInfo: fileInfo,
	}

	m.inputs[0] = textinput.New()
	// m.inputs[0].Placeholder = mdf
	m.inputs[0].Width = 50
	m.inputs[0].SetValue(fileInfo.mdfName)
	m.inputs[0].Focus()

	return m
}

func (m restoreModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m restoreModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return goToActionMsg{} }

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// check for enter on last input
			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				return m, func() tea.Msg {
					return restoreExecMsg{
						m.bakFileInfo,
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
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
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
	subHeader := TableTitleStyle.Render(fmt.Sprintf("Contents of backup file %s", m.bakFileInfo.filename))

	rowHeader := fmt.Sprintf("%s%s%s%s", ColHeaderStyle.Width(30).Render("Filename"), ColHeaderStyle.Width(10).Render("Type"), ColHeaderStyle.Width(15).Render("Size"), ColHeaderStyle.Width(15).Render("BackupSize"))
	mdfLine := fmt.Sprintf("%-30s%-10s%-15s%-15s", m.bakFileInfo.mdfName, "MDF", m.bakFileInfo.mdfSize, m.bakFileInfo.mdfBackupSize)
	ldfLine := fmt.Sprintf("%-30s%-10s%-15s%-15s", m.bakFileInfo.ldfName, "LDF", m.bakFileInfo.ldfSize, "-")

	contents := BorderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			rowHeader,
			mdfLine,
			ldfLine,
		),
	)

	input := fmt.Sprintf("New database name: %s", m.inputs[0].View())

	return lipgloss.JoinVertical(lipgloss.Left,
		subHeader,
		contents,
		input,
		"\nPress 'Enter' to confirm restore. Press 'Esc' to cancel.",
	)
}
