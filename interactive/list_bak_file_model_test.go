package interactive

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// newTestListFilesModel constructs a listFilesModel from pre-built filenames
// without calling mssql.ListBackupFilesInContainer.
func newTestListFilesModel(filenames []string) listFilesModel {
	columns := []table.Column{
		{Title: "Filename", Width: 40},
		{Title: "Size", Width: 15},
		{Title: "Date", Width: 20},
	}
	rows := make([]table.Row, len(filenames))
	for i, name := range filenames {
		rows[i] = table.Row{name, "10 MB", "2024-06-01 12:00"}
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	return listFilesModel{table: t, help: help.New()}
}

// ── View tests ────────────────────────────────────────────────────────────────

func TestListFilesModel_View_ContainsFilenames(t *testing.T) {
	m := newTestListFilesModel([]string{"vej-eu.bak", "mydb.bak"})
	view := m.View()

	if !strings.Contains(view, "vej-eu.bak") {
		t.Errorf("list-files view should contain 'vej-eu.bak'; got:\n%s", view)
	}
	if !strings.Contains(view, "mydb.bak") {
		t.Errorf("list-files view should contain 'mydb.bak'; got:\n%s", view)
	}
}

func TestListFilesModel_View_ErrorState(t *testing.T) {
	m := listFilesModel{err: &mockError{"permission denied"}}
	view := m.View()

	if !strings.Contains(view, "permission denied") {
		t.Errorf("list-files error view should contain error message; got:\n%s", view)
	}
}

// ── Update tests ──────────────────────────────────────────────────────────────

func TestListFilesModel_Update_EnterEmitsBakFileSelectedMsg(t *testing.T) {
	m := newTestListFilesModel([]string{"vej-eu.bak", "mydb.bak"})
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected a cmd from Enter, got nil")
	}
	msg := cmd()
	sel, ok := msg.(bakFileSelectedMsg)
	if !ok {
		t.Fatalf("expected bakFileSelectedMsg, got %T", msg)
	}
	if sel.filename != "vej-eu.bak" {
		t.Errorf("filename: got %q, want %q", sel.filename, "vej-eu.bak")
	}
}

func TestListFilesModel_Update_EscEmitsGoToActionMsg(t *testing.T) {
	m := newTestListFilesModel([]string{"vej-eu.bak"})
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Fatal("expected a cmd from Esc, got nil")
	}
	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Errorf("expected goToActionMsg, got %T", msg)
	}
}

func TestListFilesModel_Update_QEmitsGoToActionMsg(t *testing.T) {
	m := newTestListFilesModel([]string{"vej-eu.bak"})
	_, cmd := m.Update(tea.KeyMsg{Runes: []rune{'q'}, Type: tea.KeyRunes})

	if cmd == nil {
		t.Fatal("expected a cmd from 'q', got nil")
	}
	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Errorf("expected goToActionMsg from 'q', got %T", msg)
	}
}
