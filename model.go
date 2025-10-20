package main

import (
	"github.com/Eventid3/sql-archiver/steps"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepConnection step = iota
	stepSelectBackupFile
	stepShowDatabases
	stepRestoreOptions
	stepConfirm
	stepExecuting
)

type model struct {
	state         step
	width, height int

	// Step models
	connection steps.ConnectionModel
	filePicker steps.FilePickerModel
	// dbList        *steps.DatabaseListModel
	// restoreOpts   steps.RestoreOptionsModel
	// confirm       steps.ConfirmModel
	// executing     steps.ExecutingModel

	// Shared data
	selectedFile string
	databases    []string
	serverConfig ServerConfig
	err          error
}

type ServerConfig struct {
	Host     string
	Username string
	Password string
}

func InitialModel() model {
	return model{
		state:      stepConnection,
		connection: steps.NewConnectionModel(),
	}
}

func (m model) Init() tea.Cmd {
	return m.connection.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Handle state transitions and delegate to current step
	return m.handleStateUpdate(msg)
}

func (m model) View() string {
	switch m.state {
	case stepConnection:
		return m.connection.View()
		// case stepSelectBackupFile:
		//     return m.filePicker.View()
		// case stepShowDatabases:
		//     if m.dbList != nil {
		//         return m.dbList.View()
		//     }
		//     return "Loading databases..."
		// ... other cases
	}
	return ""
}

// handleStateUpdate manages step transitions
func (m model) handleStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConnectionEstablishedMsg:
		m.serverConfig = msg.Config
		m.state = stepSelectBackupFile
		m.filePicker = steps.NewFilePickerModel()
		return m, m.filePicker.Init()

		// case FileSelectedMsg:
		//     m.selectedFile = msg.Path
		//     m.state = stepShowDatabases
		//     m.dbList = steps.NewDatabaseListModel()
		//     return m, readDatabasesCmd(msg.Path)

		// ... other transitions
	}

	// Delegate to current step
	switch m.state {
	case stepConnection:
		newModel, cmd := m.connection.Update(msg)
		m.connection = newModel
		return m, cmd
		// case stepSelectBackupFile:
		// 	newModel, cmd := m.filePicker.Update(msg)
		// 	m.filePicker = newModel
		// 	return m, cmd
		// ... other delegations
	}

	return m, nil
}
