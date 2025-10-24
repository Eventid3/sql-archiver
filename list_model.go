package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ------ WELCOME MODEL --------
type listModel struct {
	databases []dbItem
}

type dbItem struct {
	name           string
	id             int
	created, state string
}

func NewListModel() listModel {
	return listModel{databases: getDatabaseList()}
}

func getDatabaseList() []dbItem {
	return []dbItem{
		{"db1", 1, "2023-01-01", "online"},
		{"db2", 2, "2023-01-01", "online"},
		{"db2", 3, "2023-01-01", "online"},
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return goToActionMsg{} }
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m listModel) View() string {
	header := fmt.Sprintf("\n%-30s %-12s %-20s %-15s\n", "DATABASE NAME", "ID", "CREATED", "STATE\n"+strings.Repeat("-", 80)+"\n")
	dbList := ""
	for _, db := range m.databases {
		dbList += fmt.Sprintf("%-30s %-12d %-20s %-15s\n", db.name, db.id, db.created, db.state)
	}
	return fmt.Sprintf("Here's a list of the databases:\n\n%s%s", header, dbList)
}
