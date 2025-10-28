/*
Package interactive implements an interactive command-line interface using the Bubble Tea framework.
*/
package interactive

import (
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepLogin step = iota
	stepAction
	stepList
	stepBackupSelect
	stepBackupExec
	stepRestoreSelect
	stepRestoreExec
	stepListBakFiles
	stepInspectBakFile
)

type ServerConfig struct {
	container string
	user      string
	password  string
}

type model struct {
	state      step
	form       formModel
	action     actionModel
	list       listModel
	backup     backupModel
	backupExec backupExecModel
	listFiles  listFilesModel

	activeModel tea.Model

	serverConfig ServerConfig
}

func InitialModel() model {
	form := NewFormModel()
	return model{
		state: stepLogin,
		form:  form,
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
		m.state = stepAction
		m.serverConfig.container = msg.container
		m.serverConfig.user = msg.user
		m.serverConfig.password = msg.password
		m.action = NewActionModel()
		return m, m.action.Init()
	case actionSelectedMsg:
		switch msg.action {
		case "list":
			m.state = stepList
			m.list = NewListModel(m.serverConfig.container, m.serverConfig.user, m.serverConfig.password)
			return m, m.list.Init()
		case "backup":
			m.state = stepBackupSelect
			m.backup = NewBackupModel(m.serverConfig.container, m.serverConfig.user, m.serverConfig.password)
			return m, m.backup.Init()
		case "list_files":
			m.state = stepListBakFiles
			m.listFiles = NewListFilesModel(m.serverConfig.container, m.serverConfig.user, m.serverConfig.password)
			return m, m.listFiles.Init()
		}
	case goToActionMsg:
		m.state = stepAction
		m.action = NewActionModel()
		return m, m.action.Init()
	case dbSelectedMsg:
		m.state = stepBackupExec
		m.backupExec = NewBackupExecModel(m.serverConfig.container, m.serverConfig.user, m.serverConfig.password, msg.db, msg.filename)

	}

	// HANDLE STATE UPDATES
	// -------------------
	switch m.state {
	case stepLogin:
		newForm, cmd := m.form.Update(msg)
		m.form = newForm.(formModel)
		return m, cmd
	case stepAction:
		newConfirm, cmd := m.action.Update(msg)
		m.action = newConfirm.(actionModel)
		return m, cmd
	case stepList:
		newList, cmd := m.list.Update(msg)
		m.list = newList.(listModel)
		return m, cmd
	case stepBackupSelect:
		newBackup, cmd := m.backup.Update(msg)
		m.backup = newBackup.(backupModel)
		return m, cmd
	case stepBackupExec:
		newBackupExec, cmd := m.backupExec.Update(msg)
		m.backupExec = newBackupExec.(backupExecModel)
		return m, cmd
	case stepListBakFiles:
		newListBakFiles, cmd := m.listFiles.Update(msg)
		m.listFiles = newListBakFiles.(listFilesModel)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stepLogin:
		return m.form.View()
	case stepAction:
		return m.action.View()
	case stepList:
		return m.list.View()
	case stepBackupSelect:
		return m.backup.View()
	case stepBackupExec:
		return m.backupExec.View()
	case stepListBakFiles:
		return m.listFiles.View()
	}
	return ""
}
