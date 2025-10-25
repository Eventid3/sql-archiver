package interactive

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ----------- CONFIRM MODEL -------------
type actionModel struct {
	list list.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc, action string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewActionModel() actionModel {
	items := []list.Item{
		item{title: "List databases", desc: "Show all databases in the server", action: "list"},
		item{title: "Backup", desc: "Backup selected databases to .bak file", action: "backup"},
		item{title: "Restore", desc: "Restore databases from .bak file", action: "restore"},
		item{title: "List files", desc: "List all the .bak files in the docker container", action: "list_files"},
	}

	m := actionModel{
		list: list.New(items, list.NewDefaultDelegate(), 40, 20),
	}
	m.list.Title = "Select an action"
	return m
}

// Init implements tea.Model.
func (m actionModel) Init() tea.Cmd {
	return nil
}

func (m actionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			action := m.list.SelectedItem().(item).action
			return m, func() tea.Msg { return actionSelectedMsg{action: action} }
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m actionModel) View() string {
	return docStyle.Render(m.list.View())
}
