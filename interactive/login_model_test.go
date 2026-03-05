package interactive

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestLoginModel_View_ContainsHostField(t *testing.T) {
	m := NewLoginModel(nil)
	view := m.View()

	if !strings.Contains(view, "Host:") {
		t.Errorf("login view should contain 'Host:' field label; got:\n%s", view)
	}
	if !strings.Contains(view, "User:") {
		t.Errorf("login view should contain 'User:' field label; got:\n%s", view)
	}
	if !strings.Contains(view, "Password:") {
		t.Errorf("login view should contain 'Password:' field label; got:\n%s", view)
	}
}

func TestLoginModel_View_ShowsErrorWhenPresent(t *testing.T) {
	err := &mockError{"connection refused"}
	m := NewLoginModel(err)
	view := m.View()

	if !strings.Contains(view, "connection refused") {
		t.Errorf("login view should display the error; got:\n%s", view)
	}
}

func TestLoginModel_Update_TabAdvancesFocus(t *testing.T) {
	m := NewLoginModel(nil)

	// Initial focus is index 0 (host field)
	if m.focusIndex != 0 {
		t.Fatalf("expected initial focusIndex=0, got %d", m.focusIndex)
	}

	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ := m.Update(tabMsg)
	lm := newModel.(loginModel)

	if lm.focusIndex != 1 {
		t.Errorf("after Tab, expected focusIndex=1, got %d", lm.focusIndex)
	}
}

func TestLoginModel_Update_EscQuitsApp(t *testing.T) {
	m := NewLoginModel(nil)
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := m.Update(escMsg)

	if cmd == nil {
		t.Fatal("expected a cmd from pressing Esc, got nil")
	}
	// tea.Quit returns tea.QuitMsg when executed
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg from Esc, got %T", msg)
	}
}

type mockError struct{ msg string }

func (e *mockError) Error() string { return e.msg }
