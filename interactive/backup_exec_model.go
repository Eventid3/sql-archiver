package interactive

import (
	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
)

type backupExecModel struct {
	err error
}

func NewBackupExecModel(config ServerConfig, dbName, filename string) backupExecModel {
	err := mssql.BackupDatabase(config.container, config.user, config.password, dbName, filename)

	return backupExecModel{
		err: err,
	}
}

func (m backupExecModel) Init() tea.Cmd {
	return nil
}

func (m backupExecModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return goToActionMsg{} }
		}
	}
	return m, nil
}

func (m backupExecModel) View() string {
	if m.err != nil {
		return "Error during backup:\n\n" + m.err.Error() + "\n\nPress Enter to go back to action selection."
	}

	return "Backup completed successfully!\n\nPress Enter to go back to action selection."
}
