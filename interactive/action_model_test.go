package interactive

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestActionModel_View_ContainsBothActions(t *testing.T) {
	m := NewActionModel()
	view := m.View()

	// The list renders the item titles as defined in NewActionModel
	if !strings.Contains(view, "BACKUP") {
		t.Errorf("action view should contain 'BACKUP' option; got:\n%s", view)
	}
	if !strings.Contains(view, "RESTORE") {
		t.Errorf("action view should contain 'RESTORE' option; got:\n%s", view)
	}
}

func TestActionModel_Update_EnterEmitsActionSelectedMsg(t *testing.T) {
	m := NewActionModel()

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.Update(enterMsg)

	if cmd == nil {
		t.Fatal("expected a cmd from pressing Enter, got nil")
	}

	msg := cmd()
	asMsg, ok := msg.(actionSelectedMsg)
	if !ok {
		t.Fatalf("expected actionSelectedMsg, got %T", msg)
	}
	// The default selected item should be "backup" or "restore" — either is valid
	if asMsg.action != "backup" && asMsg.action != "restore" {
		t.Errorf("unexpected action %q, want 'backup' or 'restore'", asMsg.action)
	}
}

func TestActionModel_Update_EscQuitsApp(t *testing.T) {
	m := NewActionModel()
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := m.Update(escMsg)

	if cmd == nil {
		t.Fatal("expected a cmd from pressing Esc, got nil")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg from Esc, got %T", cmd())
	}
}
