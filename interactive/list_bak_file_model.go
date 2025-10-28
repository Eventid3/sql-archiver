package interactive

import (
	"fmt"

	"github.com/Eventid3/sql-archiver/mssql"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listFilesModel struct {
	list list.Model
	err  error
}

type bakFileItem struct {
	size, date, name string
}

func (b bakFileItem) Title() string       { return b.name }
func (b bakFileItem) Description() string { return fmt.Sprintf("Size: %s, Date: %s", b.size, b.date) }
func (b bakFileItem) FilterValue() string { return b.name }

func NewListFilesModel(container, user, password string) listFilesModel {
	files, err := mssql.ListBackupFilesInContainer(container, user, password)
	if err != nil {
		return listFilesModel{
			err: err,
		}
	}

	items := []list.Item{}

	for _, f := range files {
		items = append(items, bakFileItem{
			size: f.Size,
			date: f.Date,
			name: f.Name,
		})
	}
	return listFilesModel{
		list: list.New(items, list.NewDefaultDelegate(), 40, 40),
		err:  err,
	}
}

func (m listFilesModel) Init() tea.Cmd {
	return nil
}

func (m listFilesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			file := m.list.SelectedItem().(bakFileItem)
			return m, func() tea.Msg { return bakFileSelectedMsg{file.name} }
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listFilesModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("error listing bak files: %v", m.err)
	}
	return docStyle.Render(m.list.View())
}
