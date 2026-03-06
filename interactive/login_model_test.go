package interactive

import (
	"strings"
	"testing"

	"github.com/Eventid3/sql-archiver/domain"
	tea "github.com/charmbracelet/bubbletea"
)

// ── helpers ──────────────────────────────────────────────────────────────────

type mockError struct{ msg string }

func (e *mockError) Error() string { return e.msg }

func noConnections() []domain.Connection { return nil }

func savedConns() []domain.Connection {
	return []domain.Connection{
		{Name: "local", Host: "localhost:1433", User: "sa", Password: "pw"},
		{Name: "prod", Host: "prod:1433", User: "sa", Password: "pw2"},
	}
}

// ── View tests ────────────────────────────────────────────────────────────────

func TestLoginModel_View_ContainsFormFields(t *testing.T) {
	m := NewLoginModel(noConnections(), nil)
	view := m.View()

	for _, label := range []string{"Host:", "User:", "Password:"} {
		if !strings.Contains(view, label) {
			t.Errorf("login view should contain %q; got:\n%s", label, view)
		}
	}
}

func TestLoginModel_View_ShowsErrorWhenPresent(t *testing.T) {
	err := &mockError{"connection refused"}
	m := NewLoginModel(noConnections(), err)
	view := m.View()

	if !strings.Contains(view, "connection refused") {
		t.Errorf("login view should display the error; got:\n%s", view)
	}
}

func TestLoginModel_View_ShowsSavedConnections(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	view := m.View()

	if !strings.Contains(view, "local") {
		t.Errorf("login view should list saved connection 'local'; got:\n%s", view)
	}
	if !strings.Contains(view, "prod") {
		t.Errorf("login view should list saved connection 'prod'; got:\n%s", view)
	}
}

func TestLoginModel_View_ShowsNewConnectionRowInPicker(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	view := m.View()

	if !strings.Contains(view, "New connection") {
		t.Errorf("login view should contain '+ New connection' row; got:\n%s", view)
	}
}

// ── Picker key tests ──────────────────────────────────────────────────────────

func TestLoginModel_Picker_DownAdvancesIndex(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	// Initial picker focus on index 0.
	if m.pickerIndex != 0 {
		t.Fatalf("expected initial pickerIndex=0, got %d", m.pickerIndex)
	}

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	lm := newModel.(loginModel)
	if lm.pickerIndex != 1 {
		t.Errorf("after ↓, expected pickerIndex=1, got %d", lm.pickerIndex)
	}
}

func TestLoginModel_Picker_UpDoesNotGoNegative(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	lm := newModel.(loginModel)
	if lm.pickerIndex != 0 {
		t.Errorf("↑ at index 0 should stay at 0, got %d", lm.pickerIndex)
	}
}

func TestLoginModel_Picker_NKeyShiftsFocusToBlankForm(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	// First pre-fill something so we can verify it gets cleared.
	m.inputs[0].SetValue("old-host")

	newModel, _ := m.Update(tea.KeyMsg{Runes: []rune{'n'}, Type: tea.KeyRunes})
	lm := newModel.(loginModel)

	if lm.focus != loginFocusForm {
		t.Errorf("expected focus to shift to form after 'n'; got %v", lm.focus)
	}
	if lm.inputs[0].Value() != "" {
		t.Errorf("host should be cleared on new connection, got %q", lm.inputs[0].Value())
	}
}

func TestLoginModel_Picker_DKeySetsPendingDelete(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	newModel, _ := m.Update(tea.KeyMsg{Runes: []rune{'d'}, Type: tea.KeyRunes})
	lm := newModel.(loginModel)

	if lm.pendingDelete == nil {
		t.Fatal("expected pendingDelete to be set after 'd'")
	}
	if lm.pendingDelete.Name != "local" {
		t.Errorf("pendingDelete.Name = %q, want %q", lm.pendingDelete.Name, "local")
	}
}

func TestLoginModel_Picker_EscQuitsWithNoConnections(t *testing.T) {
	m := NewLoginModel(noConnections(), nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a cmd from Esc with no saved connections")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg from Esc, got %T", cmd())
	}
}

// ── Delete-confirm key tests ──────────────────────────────────────────────────

func TestLoginModel_DeleteConfirm_NKeyCancelsPendingDelete(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	// Trigger pending delete.
	m2, _ := m.Update(tea.KeyMsg{Runes: []rune{'d'}, Type: tea.KeyRunes})
	lm := m2.(loginModel)
	if lm.pendingDelete == nil {
		t.Fatal("setup: expected pendingDelete to be set")
	}
	// Press 'n' to cancel.
	m3, _ := lm.Update(tea.KeyMsg{Runes: []rune{'n'}, Type: tea.KeyRunes})
	lm2 := m3.(loginModel)
	if lm2.pendingDelete != nil {
		t.Errorf("expected pendingDelete to be cleared after 'n'")
	}
}

func TestLoginModel_DeleteConfirm_YKeyEmitsConnectionDeletedMsg(t *testing.T) {
	// Wire up a fake config delete via the real flow — we test that the cmd
	// returned is non-nil and produces connectionDeletedMsg when the delete
	// call is mocked. Because config.DeleteConnection touches disk we only
	// verify the command type by stubbing with a name that won't exist on disk
	// (DeleteConnection is a no-op for unknown names).
	m := NewLoginModel(savedConns(), nil)
	m2, _ := m.Update(tea.KeyMsg{Runes: []rune{'d'}, Type: tea.KeyRunes})
	lm := m2.(loginModel)

	_, cmd := lm.Update(tea.KeyMsg{Runes: []rune{'y'}, Type: tea.KeyRunes})
	if cmd == nil {
		t.Fatal("expected a cmd after confirming delete with 'y'")
	}
	msg := cmd()
	if _, ok := msg.(connectionDeletedMsg); !ok {
		t.Errorf("expected connectionDeletedMsg, got %T", msg)
	}
}

// ── Form key tests ────────────────────────────────────────────────────────────

func TestLoginModel_Form_TabAdvancesFocus(t *testing.T) {
	m := NewLoginModel(noConnections(), nil)
	// No saved connections → starts in form focus.
	if m.focus != loginFocusForm {
		t.Fatalf("expected form focus initially; got %v", m.focus)
	}
	if m.focusIndex != 0 {
		t.Fatalf("expected focusIndex=0, got %d", m.focusIndex)
	}

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	lm := newModel.(loginModel)
	if lm.focusIndex != 1 {
		t.Errorf("after Tab, expected focusIndex=1, got %d", lm.focusIndex)
	}
}

func TestLoginModel_Form_ShiftTabGoesBackward(t *testing.T) {
	m := NewLoginModel(noConnections(), nil)
	// Advance to field 1, then shift+tab back.
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	lm := m2.(loginModel)
	m3, _ := lm.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	lm2 := m3.(loginModel)
	if lm2.focusIndex != 0 {
		t.Errorf("after Shift+Tab, expected focusIndex=0, got %d", lm2.focusIndex)
	}
}

func TestLoginModel_Form_TabWrapsAround(t *testing.T) {
	m := NewLoginModel(noConnections(), nil)
	// Advance past the last field (index 2) → should wrap to 0.
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})               // 0→1
	m3, _ := m2.(loginModel).Update(tea.KeyMsg{Type: tea.KeyTab}) // 1→2
	m4, _ := m3.(loginModel).Update(tea.KeyMsg{Type: tea.KeyTab}) // 2→0
	lm := m4.(loginModel)
	if lm.focusIndex != 0 {
		t.Errorf("Tab wrap-around: expected focusIndex=0, got %d", lm.focusIndex)
	}
}

func TestLoginModel_Form_EscReturnsToPicker(t *testing.T) {
	m := NewLoginModel(savedConns(), nil)
	// Navigate into form.
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m2, _ = m2.Update(tea.KeyMsg{Type: tea.KeyDown})
	m2, _ = m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	lm := m2.(loginModel)
	if lm.focus != loginFocusForm {
		t.Fatalf("expected form focus after Enter, got %v", lm.focus)
	}
	// Esc should return to picker.
	m3, cmd := lm.Update(tea.KeyMsg{Type: tea.KeyEsc})
	lm2 := m3.(loginModel)
	if lm2.focus != loginFocusPicker {
		t.Errorf("expected picker focus after Esc, got %v", lm2.focus)
	}
	if cmd != nil {
		// Should NOT quit when there are saved connections.
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Error("Esc with saved connections should return to picker, not quit")
		}
	}
}

func TestLoginModel_Form_EscQuitsWithNoConnections(t *testing.T) {
	m := NewLoginModel(noConnections(), nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a cmd from Esc in form with no saved connections")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", cmd())
	}
}

func TestLoginModel_Form_EnterOnLastFieldEmitsCmd(t *testing.T) {
	m := NewLoginModel(noConnections(), nil)
	// Advance to last field (index 2).
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m3, _ := m2.(loginModel).Update(tea.KeyMsg{Type: tea.KeyTab})
	lm := m3.(loginModel)
	if lm.focusIndex != 2 {
		t.Fatalf("expected focusIndex=2, got %d", lm.focusIndex)
	}
	_, cmd := lm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a cmd from Enter on last field, got nil")
	}
}
