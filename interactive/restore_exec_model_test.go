package interactive

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRestoreExecModel_View_SuccessPath(t *testing.T) {
	m := restoreExecModel{err: nil, infoString: ""}
	view := m.View()

	if !strings.Contains(view, "successfully") {
		t.Errorf("restore-exec success view should contain 'successfully'; got:\n%s", view)
	}
}

func TestRestoreExecModel_View_ErrorPathContainsMessageAndInfo(t *testing.T) {
	m := restoreExecModel{
		err:        &mockError{"restore failed"},
		infoString: "RESTORE DATABASE [mydb] FROM DISK = N'/var/opt/mssql/backup/mydb.bak'",
	}
	view := m.View()

	if !strings.Contains(view, "restore failed") {
		t.Errorf("restore-exec error view should contain error message; got:\n%s", view)
	}
	if !strings.Contains(view, "RESTORE DATABASE") {
		t.Errorf("restore-exec error view should contain infoString; got:\n%s", view)
	}
}

func TestRestoreExecModel_Update_EnterEmitsGoToActionMsg(t *testing.T) {
	m := restoreExecModel{err: nil}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected a cmd from Enter, got nil")
	}
	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Errorf("expected goToActionMsg, got %T", msg)
	}
}

func TestRestoreExecModel_Update_EnterOnErrorPathAlsoEmitsGoToActionMsg(t *testing.T) {
	m := restoreExecModel{err: &mockError{"restore failed"}}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected a cmd from Enter on error path, got nil")
	}
	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Errorf("expected goToActionMsg, got %T", msg)
	}
}
