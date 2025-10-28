package interactive

import (
	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
)

type restoreExecModel struct {
	err error
}

func NewRestoreExecModel(config ServerConfig, bakfile, dbName, mdf, ldf string) restoreExecModel {
	err := mssql.RestoreDatabase(config.container, config.user, config.password, bakfile, dbName, mdf, ldf)
	return restoreExecModel{err}
}

func (m restoreExecModel) Init() tea.Cmd {
	return nil
}

func (m restoreExecModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return goToActionMsg{} }
		}
	}
	return m, nil
}

func (m restoreExecModel) View() string {
	if m.err != nil {
		return "Error during restore:\n\n" + m.err.Error() + "\n\nPress Enter to go back to action selection."
	}

	return "Restore completed successfully!\n\nPress Enter to go back to action selection."
}
