package interactive

import (
	"strings"
	"testing"

	"github.com/Eventid3/sql-archiver/domain"
	tea "github.com/charmbracelet/bubbletea"
)

// newInspectModelFromEntry creates an inspectModel directly from a BackupEntry,
// bypassing the mssql.InspectBackupFile call. Used in tests.
func newInspectModelFromEntry(bakFileName string, entry domain.BackupEntry) inspectModel {
	return inspectModel{
		bakFileName: bakFileName,
		bakFileInfo: entry,
		err:         nil,
	}
}

func TestInspectModel_View_ShowsMdfAndLdfNames(t *testing.T) {
	entry := domain.BackupEntry{
		MdfFile: domain.MdfEntry{
			Name:       "VEJ-EU.Staging.Umbraco",
			Size:       "276824064",
			BackupSize: "65404928",
		},
		LdfFile: domain.LdfEntry{
			Name: "VEJ-EU.Staging.Umbraco_log",
			Size: "1954545664",
		},
	}

	m := newInspectModelFromEntry("vej-eu.bak", entry)
	view := m.View()

	if !strings.Contains(view, "VEJ-EU.Staging.Umbraco") {
		t.Errorf("View() does not contain MDF logical name; got:\n%s", view)
	}
	if !strings.Contains(view, "VEJ-EU.Staging.Umbraco_log") {
		t.Errorf("View() does not contain LDF logical name; got:\n%s", view)
	}
	if !strings.Contains(view, "vej-eu.bak") {
		t.Errorf("View() does not contain the bak filename; got:\n%s", view)
	}
}

func TestInspectModel_View_EmptyNamesShowBlank(t *testing.T) {
	entry := domain.BackupEntry{} // zero value — names will be empty strings

	m := newInspectModelFromEntry("mystery.bak", entry)
	view := m.View()

	// Should render without panic and still show the filename
	if !strings.Contains(view, "mystery.bak") {
		t.Errorf("View() does not contain the bak filename; got:\n%s", view)
	}
}

func TestInspectModel_Update_EnterEmitsRestoreBackupMsg(t *testing.T) {
	entry := domain.BackupEntry{
		Filename: "vej-eu.bak",
		MdfFile:  domain.MdfEntry{Name: "VEJ-EU.Staging.Umbraco"},
		LdfFile:  domain.LdfEntry{Name: "VEJ-EU.Staging.Umbraco_log"},
	}

	m := newInspectModelFromEntry("vej-eu.bak", entry)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.Update(enterMsg)

	if cmd == nil {
		t.Fatal("expected a cmd from pressing Enter, got nil")
	}

	// Execute the cmd to get the message it produces
	msg := cmd()
	rbMsg, ok := msg.(restoreBackupMsg)
	if !ok {
		t.Fatalf("expected restoreBackupMsg, got %T", msg)
	}
	if rbMsg.fileInfo.MdfFile.Name != "VEJ-EU.Staging.Umbraco" {
		t.Errorf("fileInfo.MdfFile.Name: got %q, want %q", rbMsg.fileInfo.MdfFile.Name, "VEJ-EU.Staging.Umbraco")
	}
	if rbMsg.fileInfo.LdfFile.Name != "VEJ-EU.Staging.Umbraco_log" {
		t.Errorf("fileInfo.LdfFile.Name: got %q, want %q", rbMsg.fileInfo.LdfFile.Name, "VEJ-EU.Staging.Umbraco_log")
	}
}

func TestInspectModel_Update_EscEmitsGoToActionMsg(t *testing.T) {
	m := newInspectModelFromEntry("vej-eu.bak", domain.BackupEntry{})

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

func TestInspectModel_View_ErrorState(t *testing.T) {
	m := inspectModel{
		bakFileName: "broken.bak",
		err:         &mockError{"access denied"},
	}
	view := m.View()

	if !strings.Contains(view, "access denied") {
		t.Errorf("inspect error view should contain error message; got:\n%s", view)
	}
}
