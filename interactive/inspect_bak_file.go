package interactive

import (
	"fmt"

	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
)

type inspectModel struct {
	bakFileName string
	bakFileInfo mssql.BackupEntry
	err         error
}

func NewInspectModel(config ServerConfig, bakFileName string) inspectModel {
	fileInfo, err := mssql.InspectBackupFile(config.container, config.user, config.password, bakFileName)
	return inspectModel{
		bakFileName: bakFileName,
		bakFileInfo: fileInfo,
		err:         err,
	}
}

func (m inspectModel) Init() tea.Cmd {
	return nil
}

func (m inspectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg {
				return restoreBackupMsg{m.bakFileName, m.bakFileInfo.MdfFile.Name, m.bakFileInfo.LdfFile.Name}
			}
		}
	}
	return m, cmd
}

func (m inspectModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error inspecting backup file: %v", m.err)
	}
	result := fmt.Sprintf("Contents of backup file %s:\n\nDatabase: %s, Size: %s, BackupSize: %s\nLdf file: %s, Ldf size: %s",
		m.bakFileName,
		m.bakFileInfo.MdfFile.Name,
		m.bakFileInfo.MdfFile.Size,
		m.bakFileInfo.MdfFile.BackupSize,
		m.bakFileInfo.LdfFile.Name,
		m.bakFileInfo.LdfFile.Size,
	)

	result += "\n\nPress enter to restore file, or esc to go back."
	return result
}
