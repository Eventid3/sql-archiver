package interactive

import (
	"github.com/Eventid3/sql-archiver/domain"
	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type restoreExecModel struct {
	err        error
	infoString string
}

func NewRestoreExecModel(config ServerConfig, bakFileInfo domain.BackupEntry, newDBName string) restoreExecModel {
	query, err := mssql.RestoreDatabase(
		config.container,
		config.user,
		config.password,
		bakFileInfo.Filename,
		newDBName,
		bakFileInfo.MdfFile.Name,
		bakFileInfo.LdfFile.Name)
	return restoreExecModel{err, query}
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
		return "Error during restore:\n\n" + m.err.Error() + "\n\nPress Enter to go back to action selection.\nInfo:\n" + m.infoString
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		"Restore completed successfully!",
		"Press 'Enter' to go back to action selection",
	)
}
