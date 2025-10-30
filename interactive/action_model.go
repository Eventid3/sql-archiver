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

type actionItem struct {
	title, desc, action string
}

func (i actionItem) Title() string       { return i.title }
func (i actionItem) Description() string { return i.desc }
func (i actionItem) FilterValue() string { return i.title }

func NewActionModel() actionModel {
	items := []list.Item{
		actionItem{title: "List databases  -> BACKUP", desc: "Show all databases in the server", action: "backup"},
		actionItem{title: "List .bak files -> RESTORE", desc: "List all the .bak files in the docker container, and restore database", action: "restore"},
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
		switch msg.String() {
		case "enter":
			action := m.list.SelectedItem().(actionItem).action
			return m, func() tea.Msg { return actionSelectedMsg{action: action} }
		case "esc":
			return m, tea.Quit
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
