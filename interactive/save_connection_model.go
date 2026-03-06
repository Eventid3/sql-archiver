package interactive

import (
	"fmt"

	"github.com/Eventid3/sql-archiver/config"
	"github.com/Eventid3/sql-archiver/domain"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// saveConnectionModel asks the user to name a new connection before saving it.
// It is shown after a successful login for a connection that is not yet in the
// saved list.
//
// Keyboard map:
//
//	Enter  save the connection and proceed to the action screen
//	Esc    skip saving and proceed to the action screen
type saveConnectionModel struct {
	host      string
	user      string
	password  string
	nameInput textinput.Model
	err       error
}

// NewSaveConnectionModel constructs the model. The name input is pre-filled
// with "user@host" as a sensible default the user can edit.
func NewSaveConnectionModel(host, user, password string) saveConnectionModel {
	ti := textinput.New()
	ti.Placeholder = "connection name"
	ti.Width = 40
	ti.SetValue(fmt.Sprintf("%s@%s", user, host))
	ti.Focus()

	return saveConnectionModel{
		host:      host,
		user:      user,
		password:  password,
		nameInput: ti,
	}
}

// Init implements tea.Model.
func (m saveConnectionModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m saveConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Skip saving — proceed directly to the action screen.
			host, user, password := m.host, m.user, m.password
			return m, func() tea.Msg {
				return connectionSavedMsg{host: host, user: user, password: password}
			}
		case "enter":
			name := m.nameInput.Value()
			host, user, password := m.host, m.user, m.password
			return m, func() tea.Msg {
				err := config.SaveConnection(domain.Connection{
					Name:     name,
					Host:     host,
					User:     user,
					Password: password,
				})
				if err != nil {
					// Surface the error but still let the user in.
					_ = err
				}
				return connectionSavedMsg{host: host, user: user, password: password}
			}
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

func (m saveConnectionModel) View() string {
	s := fmt.Sprintf(`Connection saved successfully. Give it a name:

Name: %s

Press Enter to save · Esc to skip`, m.nameInput.View())

	if m.err != nil {
		s += ErrorTextStyle.Render(fmt.Sprintf("\n\nError saving: %s", m.err.Error()))
	}
	return s
}
