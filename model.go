package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepWelcome step = iota
	stepForm
	stepConfirm
)

type model struct {
	state   step
	welcome welcomeModel
	form    formModel
	confirm confirmModel

	formData string
}

func InitialModel() model {
	return model{
		state:   0,
		welcome: NewWelcomeModel(),
	}
}

type nextStepMsg struct{}

type formDoneMsg struct {
	data string
}

func (m model) Init() tea.Cmd {
	m.welcome = NewWelcomeModel()
	return m.welcome.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case nextStepMsg:
		m.state++
		if m.state == stepForm {
			m.form = NewFormModel()
		}
		return m, m.form.Init()
	case formDoneMsg:
		m.state++
		m.formData = msg.data
		m.confirm = NewConfirmModel(m.formData)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
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
	case stepConfirm:
		newConfirm, cmd := m.confirm.Update(msg)
		m.confirm = newConfirm.(confirmModel)
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
	case stepConfirm:
		return m.confirm.View()
	}
	return ""
}

// ------ WELCOME MODEL --------
type welcomeModel struct {
	welcomeMsg string
}

func NewWelcomeModel() welcomeModel {
	return welcomeModel{"Welcome!"}
}

// Init implements tea.Model.
func (vm welcomeModel) Init() tea.Cmd {
	return nil
}

func (vm welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return vm, func() tea.Msg { return nextStepMsg{} }
		}
	}
	return vm, nil
}

// View implements tea.Model.
func (vm welcomeModel) View() string {
	return "Welcome to the application!\n\nPress Enter to continue."
}

// -------- FORM MODEL --------------
type formModel struct {
	input textinput.Model
}

func NewFormModel() formModel {
	input := textinput.New()
	input.Placeholder = "Type something..."
	input.Focus()
	return formModel{input}
}

// Init implements tea.Model.
func (fm formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (fm formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return fm, func() tea.Msg { return formDoneMsg{fm.input.Value()} }
		}
	}

	fm.input, cmd = fm.input.Update(msg)
	return fm, cmd
}

// View implements tea.Model.
func (fm formModel) View() string {
	return fmt.Sprintf(
		"Fill the form below:\n\n%s\n\n%s",
		fm.input.View(),
		"Enter to submit, Esc or ctrl+q to quit.",
	)
}

// ----------- CONFIRM MODEL -------------
type confirmModel struct {
	data string
}

func NewConfirmModel(msg string) confirmModel {
	return confirmModel{msg}
}

// Init implements tea.Model.
func (cm confirmModel) Init() tea.Cmd {
	return nil
}

func (cm confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch masg := msg.(type) {
	case tea.KeyMsg:
		switch masg.String() {
		case "enter":
			return cm, tea.Quit
		}
	}
	return cm, nil
}

// View implements tea.Model.
func (cm confirmModel) View() string {
	return fmt.Sprintf(
		"You entered %s!\n\n",
		cm.data,
	)
}
