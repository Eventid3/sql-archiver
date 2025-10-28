package interactive

import (
	"fmt"
	"strings"

	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
)

// ------ WELCOME MODEL --------
type listModel struct {
	databases []mssql.DBItem
	err       error
}

func NewListModel(config ServerConfig) listModel {
	dbList, err := mssql.GetDatabases(config.container, config.user, config.password)
	return listModel{databases: dbList, err: err}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return goToActionMsg{} }
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m listModel) View() string {
	if m.err != nil {
		return fmt.Sprintf(
			"Error retrieving databases:\n\n%v\n\n%s",
			m.err,
			"Press Enter to go back to action selection.",
		)
	} else {
		header := fmt.Sprintf("\n%-30s %-12s %-20s %-15s\n", "DATABASE NAME", "ID", "CREATED", "STATE\n"+strings.Repeat("-", 80))
		dbList := ""
		for _, item := range m.databases {
			if len(item.Name) > 28 {
				item.Name = item.Name[:25] + "..."
			}
			dbList += fmt.Sprintf("%-30s %-12s %-20s %-15s\n", item.Name, item.ID, item.Created, item.State)
		}
		return fmt.Sprintf(
			"Here's a list of the databases:\n\n%s%s\n\n%s",
			header,
			dbList,
			"Press Enter to go back to action selection.",
		)
	}
}
