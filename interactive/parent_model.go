/*
Package interactive implements an interactive command-line interface using the Bubble Tea framework.
*/
package interactive

import (
	"github.com/Eventid3/sql-archiver/config"
	"github.com/Eventid3/sql-archiver/domain"
	"github.com/Eventid3/sql-archiver/mssql"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ServerConfig holds the active SQL Server connection credentials.
type ServerConfig struct {
	host     string
	user     string
	password string
}

type parentModel struct {
	activeModel  tea.Model
	serverConfig ServerConfig
	connections  []domain.Connection // saved connections from disk
}

// InitialModel is the entry point when no config file is available.
// It loads any persisted connections from disk and shows the login screen.
func InitialModel() *parentModel {
	conns, _ := config.LoadConnections()
	return &parentModel{
		activeModel: NewLoginModel(conns, nil),
		connections: conns,
	}
}

// InitialModelWithConfig is the entry point when the config file provides a
// single pre-loaded connection (legacy path kept for backward compatibility).
// The connection is attempted immediately; on failure the login screen is shown.
func InitialModelWithConfig(host, user, password string) *parentModel {
	conns, _ := config.LoadConnections()
	cfg := ServerConfig{host, user, password}

	if err := mssql.CheckConnection(host, user, password); err != nil {
		return &parentModel{
			activeModel: NewLoginModel(conns, err),
			connections: conns,
		}
	}

	return &parentModel{
		activeModel:  NewActionModel(),
		serverConfig: cfg,
		connections:  conns,
	}
}

// Init implements tea.Model.
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
		m.activeModel = NewLoginModel(m.connections, msg.err)
		return m, m.activeModel.Init()

	case loginDoneMsg:
		m.serverConfig.host = msg.host
		m.serverConfig.user = msg.user
		m.serverConfig.password = msg.password
		m.activeModel = NewActionModel()
		return m, m.activeModel.Init()

	case promptSaveConnectionMsg:
		m.serverConfig.host = msg.host
		m.serverConfig.user = msg.user
		m.serverConfig.password = msg.password
		m.activeModel = NewSaveConnectionModel(msg.host, msg.user, msg.password)
		return m, m.activeModel.Init()

	case connectionSavedMsg:
		// Refresh the saved-connections list from disk.
		conns, _ := config.LoadConnections()
		m.connections = conns
		m.activeModel = NewActionModel()
		return m, m.activeModel.Init()

	case connectionDeletedMsg:
		// Reload connections and show the login screen with the updated list.
		conns, _ := config.LoadConnections()
		m.connections = conns
		m.activeModel = NewLoginModel(conns, nil)
		return m, m.activeModel.Init()

	case goToActionMsg:
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

	case dbSelectedMsg:
		m.activeModel = NewBackupExecModel(m.serverConfig, msg.db, msg.filename)
		return m, m.activeModel.Init()

	case bakFileSelectedMsg:
		m.activeModel = NewInspectModel(m.serverConfig, msg.filename)
		return m, m.activeModel.Init()

	case restoreBackupMsg:
		m.activeModel = NewRestoreModel(m.serverConfig, msg.fileInfo)
		return m, m.activeModel.Init()

	case restoreExecMsg:
		m.activeModel = NewRestoreExecModel(m.serverConfig, msg.fileInfo, msg.newDBName)
		return m, m.activeModel.Init()
	}

	// HANDLE STATE UPDATES
	// -------------------
	newActiveModel, cmd := m.activeModel.Update(msg)
	m.activeModel = newActiveModel
	return m, cmd
}

func (m *parentModel) View() string {
	return OuterStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		HeadingStyle.Render(Logo),
		RenderStatusBar(m.serverConfig.host, m.activeModel),
		"\n",
		m.activeModel.View(),
	))
}
