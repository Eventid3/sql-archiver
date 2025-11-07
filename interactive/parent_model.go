/*
Package interactive implements an interactive command-line interface using the Bubble Tea framework.
*/
package interactive

import (
	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ServerConfig struct {
	container string
	user      string
	password  string
}

type parentModel struct {
	activeModel  tea.Model
	serverConfig ServerConfig
}

func InitialModel() *parentModel {
	form := NewLoginModel(nil)
	return &parentModel{
		activeModel: form,
	}
}

func InitialModelWithConfig(container, user, password string) *parentModel {
	config := ServerConfig{
		container,
		user,
		password,
	}

	err := mssql.CheckConnection(container, user, password)
	if err != nil {
		return &parentModel{
			activeModel: NewLoginModel(err),
		}
	}

	return &parentModel{
		activeModel:  NewActionModel(),
		serverConfig: config,
	}
}

func (m *parentModel) Init() tea.Cmd {
	return m.activeModel.Init()
}

func (m *parentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// INTERCEPT MESSAGES
	// -------------------
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case loginFailedMsg:
		m.activeModel = NewLoginModel(msg.err)
		return m, m.activeModel.Init()
	case loginDoneMsg:
		m.serverConfig.container = msg.container
		m.serverConfig.user = msg.user
		m.serverConfig.password = msg.password
		m.activeModel = NewActionModel()
		return m, m.activeModel.Init()
	case actionSelectedMsg:
		switch msg.action {
		case "backup":
			m.activeModel = NewBackupModel(m.serverConfig)
			return m, m.activeModel.Init()
		case "restore":
			m.activeModel = NewListFilesModel(m.serverConfig)
			return m, m.activeModel.Init()
		}
	case goToActionMsg:
		m.activeModel = NewActionModel()
		return m, m.activeModel.Init()
	case dbSelectedMsg:
		m.activeModel = NewBackupExecModel(m.serverConfig, msg.db, msg.filename)
		return m, m.activeModel.Init()
	case bakFileSelectedMsg:
		m.activeModel = NewInspectModel(m.serverConfig, msg.filename)
		return m, m.activeModel.Init()
	case restoreBackupMsg:
		m.activeModel = NewRestoreModel(m.serverConfig, msg.filename, msg.mdfName, msg.ldfName)
		return m, m.activeModel.Init()
	case restoreExecMsg:
		m.activeModel = NewRestoreExecModel(m.serverConfig, msg.filename, msg.newDBName, msg.mdfName, msg.ldfName)
		return m, m.activeModel.Init()
	}

	// HANDLE STATE UPDATES
	// -------------------
	newActiveModel, cmd := m.activeModel.Update(msg)
	m.activeModel = newActiveModel
	return m, cmd
}

func (m *parentModel) View() string {
	header := lipgloss.JoinVertical(lipgloss.Left,
		HeadingStyle.Render(Logo),
		RenderStatusBar(m.serverConfig.container, m.activeModel),
	)

	return OuterStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		header, "\n",
		m.activeModel.View(),
	))
}
