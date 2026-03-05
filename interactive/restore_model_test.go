package interactive

import (
	"strings"
	"testing"

	"github.com/Eventid3/sql-archiver/domain"
	tea "github.com/charmbracelet/bubbletea"
)

func TestRestoreModel_View_PreFillsDBNameFromMdf(t *testing.T) {
	entry := domain.BackupEntry{
		Filename: "vej-eu.bak",
		MdfFile:  domain.MdfEntry{Name: "VEJ-EU.Staging.Umbraco"},
		LdfFile:  domain.LdfEntry{Name: "VEJ-EU.Staging.Umbraco_log"},
	}

	m := NewRestoreModel(ServerConfig{}, entry)
	view := m.View()

	// The DB name input should be pre-filled with the MDF logical name
	if !strings.Contains(view, "VEJ-EU.Staging.Umbraco") {
		t.Errorf("restore view should pre-fill with MDF logical name; got:\n%s", view)
	}
}

func TestRestoreModel_Update_EnterEmitsRestoreExecMsg(t *testing.T) {
	entry := domain.BackupEntry{
		Filename: "vej-eu.bak",
		MdfFile:  domain.MdfEntry{Name: "VEJ-EU.Staging.Umbraco"},
		LdfFile:  domain.LdfEntry{Name: "VEJ-EU.Staging.Umbraco_log"},
	}

	m := NewRestoreModel(ServerConfig{}, entry)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.Update(enterMsg)

	if cmd == nil {
		t.Fatal("expected a cmd from pressing Enter, got nil")
	}

	msg := cmd()
	reMsg, ok := msg.(restoreExecMsg)
	if !ok {
		t.Fatalf("expected restoreExecMsg, got %T", msg)
	}
	if reMsg.fileInfo.Filename != "vej-eu.bak" {
		t.Errorf("fileInfo.Filename: got %q, want %q", reMsg.fileInfo.Filename, "vej-eu.bak")
	}
	if reMsg.fileInfo.MdfFile.Name != "VEJ-EU.Staging.Umbraco" {
		t.Errorf("fileInfo.MdfFile.Name: got %q, want %q", reMsg.fileInfo.MdfFile.Name, "VEJ-EU.Staging.Umbraco")
	}
}

func TestRestoreModel_Update_EscEmitsGoToActionMsg(t *testing.T) {
	m := NewRestoreModel(ServerConfig{}, domain.BackupEntry{})

	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := m.Update(escMsg)

	if cmd == nil {
		t.Fatal("expected a cmd from pressing Esc, got nil")
	}

	msg := cmd()
	if _, ok := msg.(goToActionMsg); !ok {
		t.Fatalf("expected goToActionMsg, got %T", msg)
	}
}
