/*
Package interactive implements an interactive command-line interface using the Bubble Tea framework.
*/
package interactive

import (
	tea "github.com/charmbracelet/bubbletea"
)

// type step int
//
// const (
// 	stepLogin step = iota
// 	stepAction
// 	stepList
// 	stepBackupSelect
// 	stepBackupExec
// 	stepRestoreSelect
// 	stepRestoreExec
// 	stepListBakFiles
// 	stepInspectBakFile
// )

type ServerConfig struct {
	container string
	user      string
	password  string
}

type model struct {
	// state        step
	activeModel  tea.Model
	serverConfig ServerConfig
}

func InitialModel() model {
	form := NewFormModel()
	return model{
		// state:       stepLogin,
		activeModel: form,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// INTERCEPT MESSAGES
	// -------------------
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
	case loginDoneMsg:
		// m.state = stepAction
		m.serverConfig.container = msg.container
		m.serverConfig.user = msg.user
		m.serverConfig.password = msg.password
		m.activeModel = NewActionModel()
		return m, m.activeModel.Init()
	case actionSelectedMsg:
		switch msg.action {
		case "list":
			// m.state = stepList
			m.activeModel = NewListModel(m.serverConfig)
			return m, m.activeModel.Init()
		case "backup":
			// m.state = stepBackupSelect
			m.activeModel = NewBackupModel(m.serverConfig)
			return m, m.activeModel.Init()
		case "list_files":
			// m.state = stepListBakFiles
			m.activeModel = NewListFilesModel(m.serverConfig)
			return m, m.activeModel.Init()
		}
	case goToActionMsg:
		// m.state = stepAction
		m.activeModel = NewActionModel()
		return m, m.activeModel.Init()
	case dbSelectedMsg:
		// m.state = stepBackupExec
		m.activeModel = NewBackupExecModel(m.serverConfig, msg.db, msg.filename)
		return m, m.activeModel.Init()
	case bakFileSelectedMsg:
		// m.state = stepInspectBakFile
		m.activeModel = NewInspectModel(m.serverConfig, msg.filename)
		return m, m.activeModel.Init()
	case restoreBackupMsg:
		// m.state = stepRestoreSelect
		m.activeModel = NewRestoreModel(m.serverConfig, msg.filename, msg.mdfName, msg.ldfName)
		return m, m.activeModel.Init()
	case restoreExecMsg:
		// m.state = stepRestoreExec
		m.activeModel = NewRestoreExecModel(m.serverConfig, msg.filename, msg.newDBName, msg.mdfName, msg.ldfName)
		return m, m.activeModel.Init()
	}

	// HANDLE STATE UPDATES
	// -------------------
	newActiveModel, cmd := m.activeModel.Update(msg)
	m.activeModel = newActiveModel
	return m, cmd
}

func (m model) View() string {
	return m.activeModel.View()
}
