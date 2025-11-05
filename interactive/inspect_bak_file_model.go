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
		case "esc":
			return m, func() tea.Msg { return goToActionMsg{} }
		}
	}
	return m, cmd
}

func (m inspectModel) View() string {
	if m.err != nil {
		return errorTextStyle.Render(fmt.Sprintf("Error inspecting backup file: %v", m.err))
	}

	subHeader := baseStyle.Render(fmt.Sprintf("Contents of backup file %s", m.bakFileName))

	rowHeader := fmt.Sprintf("%s%s%s%s", colHeaderStyle.Width(30).Render("Filename"), colHeaderStyle.Width(10).Render("Type"), colHeaderStyle.Width(15).Render("Size"), colHeaderStyle.Width(15).Render("BackupSize"))
	mdfLine := fmt.Sprintf("%-30s%-10s%-15s%-15s", m.bakFileInfo.MdfFile.Name, "MDF", m.bakFileInfo.MdfFile.Size, m.bakFileInfo.MdfFile.BackupSize)
	ldfLine := fmt.Sprintf("%-30s%-10s%-15s%-15s", m.bakFileInfo.LdfFile.Name, "LDF", m.bakFileInfo.LdfFile.Size, "-")

	return subHeader + "\n" + rowHeader + "\n" + mdfLine + "\n" + ldfLine
}
