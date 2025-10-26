package interactive

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ------ WELCOME MODEL --------
type templateModel struct {
	data string
}

func NewWelcomeModel() templateModel {
	return templateModel{"Hello, World"}
}

// Init implements tea.Model.
func (m templateModel) Init() tea.Cmd {
	return nil
}

func (m templateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return nextStepMsg{} }
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m templateModel) View() string {
	return fmt.Sprintf("Here's the data: %s", m.data)
}
