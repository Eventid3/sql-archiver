package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepWelcome step = iota
	stepForm
	stepAction
	stepList
)

type model struct {
	state   step
	welcome welcomeModel
	form    formModel
	action  actionModel
	list    listModel

	formData string
}

func InitialModel() model {
	return model{
		state:   0,
		welcome: NewWelcomeModel(),
	}
}

func (m model) Init() tea.Cmd {
	m.welcome = NewWelcomeModel()
	return m.welcome.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
	case nextStepMsg:
		m.state++
		if m.state == stepForm {
			m.form = NewFormModel()
		}
		return m, m.form.Init()
	case formDoneMsg:
		m.state++
		m.formData = msg.user
		m.action = NewConfirmModel(m.formData)
		return m, m.action.Init()
	case actionSelectedMsg:
		switch msg.action {
		case "list":
			m.state = stepList
			m.list = NewListModel()
			return m, m.list.Init()
		}
	}

	switch m.state {
	case stepWelcome:
		newWelcome, cmd := m.welcome.Update(msg)
		m.welcome = newWelcome.(welcomeModel)
		return m, cmd
	case stepForm:
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
	case stepWelcome:
		return m.welcome.View()
	case stepForm:
		return m.form.View()
	case stepAction:
		return m.action.View()
	case stepList:
		return m.list.View()
	}
	return ""
}
