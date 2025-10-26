package interactive

import (
	"github.com/Eventid3/sql-archiver/mssql"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// ------ WELCOME MODEL --------
type backupModel struct {
	table    table.Model
	filename textinput.Model
	err      error
}

func NewBackupModel(container, user, password string) backupModel {
	columns := []table.Column{
		{Title: "Database", Width: 30},
		{Title: "ID", Width: 12},
		{Title: "Created", Width: 20},
		{Title: "State", Width: 15},
	}

	databases, err := mssql.GetDatabases(container, user, password)
	if err != nil {
		return backupModel{err: err}
	}

	rows := []table.Row{}

	for _, item := range databases {
		rows = append(rows, table.Row{item.Name, item.ID, item.Created, item.State})
	}

	dbTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	filenameInput := textinput.New()
	filenameInput.Placeholder = "backup_filename.bak"
	filenameInput.Width = 50

	return backupModel{
		dbTable,
		filenameInput,
		nil,
	}
}

// Init implements tea.Model.
func (m backupModel) Init() tea.Cmd {
	return nil
}

func (m backupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m backupModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
