package interactive

import (
	"fmt"

	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
)

type inspectModel struct {
	bakFileName string
	bakFileInfo string
	err         error
}

func NewInspectModel(container, user, password, bakFileName string) inspectModel {
	fileInfo, err := mssql.InspectBackupFile(container, user, password, bakFileName)
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
			return m, func() tea.Msg { return goToActionMsg{} }
		}
	}
	return m, cmd
}

func (m inspectModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error inspecting backup file: %v", m.err)
	}
	return fmt.Sprintf("Contents of backup file %s:\n\n%s", m.bakFileName, m.bakFileInfo)
}
