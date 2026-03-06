package interactive

import (
	"strings"
	"testing"

	"github.com/Eventid3/sql-archiver/domain"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// newTestBackupModel builds a backupModel from a slice of DBItems without
// touching the database.
func newTestBackupModel(items []domain.DBItem) backupModel {
	columns := []table.Column{
		{Title: "Database", Width: 30},
		{Title: "ID", Width: 12},
		{Title: "Created", Width: 20},
		{Title: "State", Width: 15},
	}
	rows := make([]table.Row, len(items))
	for i, item := range items {
		rows[i] = table.Row{item.Name, item.ID, item.Created, item.State}
	}
	dbTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	fi := textinput.New()
	fi.Placeholder = "backup_filename.bak"
	fi.Width = 50
	return backupModel{table: dbTable, filename: fi, help: help.New()}
}

// ── View tests ────────────────────────────────────────────────────────────────

func TestBackupModel_View_ContainsTableAndFilenamePrompt(t *testing.T) {
	items := []domain.DBItem{
		{Name: "mydb", ID: "5", Created: "2024-06-01", State: "ONLINE"},
	}
	m := newTestBackupModel(items)
	view := m.View()

	if !strings.Contains(view, "mydb") {
		t.Errorf("backup view should contain database name 'mydb'; got:\n%s", view)
	}
	if !strings.Contains(view, ".bak") {
		t.Errorf("backup view should contain '.bak' filename hint; got:\n%s", view)
	}
}

// ── Update tests ──────────────────────────────────────────────────────────────

func TestBackupModel_Update_EnterOnTableBlurstTableAndFocusesFilename(t *testing.T) {
	items := []domain.DBItem{
		{Name: "master", ID: "1", Created: "2024-01-01", State: "ONLINE"},
	}
	m := newTestBackupModel(items)

	if !m.table.Focused() {
		t.Fatal("expected table to be focused initially")
	}

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	bm := newModel.(backupModel)

	if bm.table.Focused() {
		t.Error("table should be blurred after Enter")
	}
}

func TestBackupModel_Update_EscOnTableEmitsGoToActionMsg(t *testing.T) {
	items := []domain.DBItem{
		{Name: "master", ID: "1", Created: "2024-01-01", State: "ONLINE"},
	}
	m := newTestBackupModel(items)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a cmd from Esc, got nil")
	}
	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Errorf("expected goToActionMsg, got %T", msg)
	}
}

func TestBackupModel_Update_EnterOnFilenameEmitsDbSelectedMsg(t *testing.T) {
	items := []domain.DBItem{
		{Name: "master", ID: "1", Created: "2024-01-01", State: "ONLINE"},
	}
	m := newTestBackupModel(items)

	// First Enter: blur table, focus filename.
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	bm := m2.(backupModel)
	bm.filename.SetValue("backup.bak")

	// Second Enter: emit dbSelectedMsg.
	_, cmd := bm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a cmd from Enter on filename, got nil")
	}
	msg := cmd()
	sel, ok := msg.(dbSelectedMsg)
	if !ok {
		t.Fatalf("expected dbSelectedMsg, got %T", msg)
	}
	if sel.db != "master" {
		t.Errorf("db: got %q, want %q", sel.db, "master")
	}
	if sel.filename != "backup.bak" {
		t.Errorf("filename: got %q, want %q", sel.filename, "backup.bak")
	}
}

func TestBackupModel_Update_EscOnFilenameRefocusesTable(t *testing.T) {
	items := []domain.DBItem{
		{Name: "master", ID: "1", Created: "2024-01-01", State: "ONLINE"},
	}
	m := newTestBackupModel(items)

	// Enter: blur table, focus filename.
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	bm := m2.(backupModel)
	if bm.table.Focused() {
		t.Fatal("setup: expected table to be blurred after Enter")
	}

	// Esc: re-focus table.
	m3, _ := bm.Update(tea.KeyMsg{Type: tea.KeyEsc})
	bm2 := m3.(backupModel)
	if !bm2.table.Focused() {
		t.Error("expected table to be focused again after Esc on filename")
	}
}
