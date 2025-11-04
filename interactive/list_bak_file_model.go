package interactive

import (
	"fmt"

	"github.com/Eventid3/sql-archiver/mssql"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type listFilesModel struct {
	table table.Model
	help  help.Model
	err   error
}

func NewListFilesModel(config ServerConfig) listFilesModel {
	files, err := mssql.ListBackupFilesInContainer(config.container, config.user, config.password)
	if err != nil {
		return listFilesModel{
			err: err,
		}
	}

	columns := []table.Column{
		{Title: "Filename", Width: 40},
		{Title: "Size", Width: 15},
		{Title: "Date", Width: 20},
	}

	rows := []table.Row{}

	for _, f := range files {
		rows = append(rows, table.Row{
			f.Name,
			f.Size,
			f.Date,
		})
	}

	fileTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	return listFilesModel{
		fileTable,
		help.New(),
		nil,
	}
}

func (m listFilesModel) Init() tea.Cmd {
	return nil
}

func (m listFilesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			filename := m.table.SelectedRow()[0]
			return m, func() tea.Msg { return bakFileSelectedMsg{filename} }
		case "esc", "q":
			return m, func() tea.Msg { return goToActionMsg{} }
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m listFilesModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("error listing bak files: %v", m.err)
	}
	return baseStyle.Render(m.table.View()) + "\n\n" +
		m.help.FullHelpView(m.table.KeyMap.FullHelp())
}
