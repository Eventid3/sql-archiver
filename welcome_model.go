package main

import tea "github.com/charmbracelet/bubbletea"

// ------ WELCOME MODEL --------
type welcomeModel struct {
	welcomeMsg string
}

func NewWelcomeModel() welcomeModel {
	return welcomeModel{"Welcome!"}
}

// Init implements tea.Model.
func (m welcomeModel) Init() tea.Cmd {
	return nil
}

func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m welcomeModel) View() string {
	return "Welcome to the application!\n\nPress Enter to continue."
}
