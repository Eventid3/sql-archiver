/*
Package interactive implements an interactive command-line interface using the Bubble Tea framework.
*/
package interactive

import (
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepLogin step = iota
	stepAction
	stepList
)

type ServerConfig struct {
	container string
	user      string
	password  string
}

type model struct {
	state  step
	form   formModel
	action actionModel
	list   listModel

	serverConfig ServerConfig
}

func InitialModel() model {
	return model{
		state: stepLogin,
		form:  NewFormModel(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
	case loginDoneMsg:
		m.state = stepAction
		m.serverConfig.container = msg.container
		m.serverConfig.user = msg.user
		m.serverConfig.password = msg.password
		m.action = NewActionModel()
		return m, m.action.Init()
	case actionSelectedMsg:
		switch msg.action {
		case "list":
			m.state = stepList
			m.list = NewListModel(m.serverConfig.container, m.serverConfig.user, m.serverConfig.password)
			return m, m.list.Init()
		}
	case goToActionMsg:
		m.state = stepAction
		m.action = NewActionModel()
		return m, m.action.Init()
	}

	switch m.state {
	case stepLogin:
		newForm, cmd := m.form.Update(msg)
		m.form = newForm.(formModel)
		return m, cmd
	case stepAction:
		newConfirm, cmd := m.action.Update(msg)
		m.action = newConfirm.(actionModel)
		return m, cmd
	case stepList:
		newList, cmd := m.list.Update(msg)
		m.list = newList.(listModel)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stepLogin:
		return m.form.View()
	case stepAction:
		return m.action.View()
	case stepList:
		return m.list.View()
	}
	return ""
}
