package interactive

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSaveConnectionModel_View_ContainsNameInput(t *testing.T) {
	m := NewSaveConnectionModel("localhost:1433", "sa", "pw")
	view := m.View()

	if !strings.Contains(view, "Name:") {
		t.Errorf("save-connection view should contain 'Name:'; got:\n%s", view)
	}
}

func TestSaveConnectionModel_View_PreFillsNameFromHostUser(t *testing.T) {
	m := NewSaveConnectionModel("localhost:1433", "sa", "pw")
	view := m.View()

	if !strings.Contains(view, "sa@localhost:1433") {
		t.Errorf("save-connection view should pre-fill name with 'sa@localhost:1433'; got:\n%s", view)
	}
}

func TestSaveConnectionModel_Update_EscSkipsSaveAndEmitsConnectionSavedMsg(t *testing.T) {
	m := NewSaveConnectionModel("localhost:1433", "sa", "pw")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Fatal("expected a cmd from Esc, got nil")
	}
	msg := cmd()
	saved, ok := msg.(connectionSavedMsg)
	if !ok {
		t.Fatalf("expected connectionSavedMsg, got %T", msg)
	}
	if saved.host != "localhost:1433" {
		t.Errorf("host: got %q, want %q", saved.host, "localhost:1433")
	}
	if saved.user != "sa" {
		t.Errorf("user: got %q, want %q", saved.user, "sa")
	}
}

func TestSaveConnectionModel_Update_EnterEmitsConnectionSavedMsg(t *testing.T) {
	m := NewSaveConnectionModel("localhost:1433", "sa", "pw")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected a cmd from Enter, got nil")
	}
	msg := cmd()
	// The cmd writes to disk which may fail in CI without a real config dir,
	// but should still return connectionSavedMsg.
	if _, ok := msg.(connectionSavedMsg); !ok {
		t.Fatalf("expected connectionSavedMsg, got %T", msg)
	}
}

func TestSaveConnectionModel_Update_CtrlCQuitsApp(t *testing.T) {
	m := NewSaveConnectionModel("localhost:1433", "sa", "pw")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected a cmd from ctrl+c, got nil")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg from ctrl+c, got %T", cmd())
	}
}
