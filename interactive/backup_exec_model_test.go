package interactive

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestBackupExecModel_View_SuccessPath(t *testing.T) {
	m := backupExecModel{err: nil}
	view := m.View()

	if !strings.Contains(view, "successfully") {
		t.Errorf("backup-exec success view should contain 'successfully'; got:\n%s", view)
	}
}

func TestBackupExecModel_View_ErrorPath(t *testing.T) {
	m := backupExecModel{err: &mockError{"disk full"}}
	view := m.View()

	if !strings.Contains(view, "disk full") {
		t.Errorf("backup-exec error view should contain error message; got:\n%s", view)
	}
}

func TestBackupExecModel_Update_EnterEmitsGoToActionMsg(t *testing.T) {
	m := backupExecModel{err: nil}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected a cmd from Enter, got nil")
	}
	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Errorf("expected goToActionMsg, got %T", msg)
	}
}

func TestBackupExecModel_Update_EnterOnErrorPathAlsoEmitsGoToActionMsg(t *testing.T) {
	m := backupExecModel{err: &mockError{"disk full"}}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected a cmd from Enter on error path, got nil")
	}
	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Errorf("expected goToActionMsg, got %T", msg)
	}
}
