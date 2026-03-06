package interactive

import (
	"testing"

	"github.com/Eventid3/sql-archiver/domain"
	tea "github.com/charmbracelet/bubbletea"
)

// newTestParent builds a parentModel in a known state without loading config
// from disk or hitting the database.
func newTestParent() *parentModel {
	return &parentModel{
		activeModel: NewLoginModel(nil, nil),
		connections: nil,
	}
}

// sendMsg dispatches a single message to the parentModel and returns the
// updated model. It mirrors what tea.Program does internally.
func sendMsg(p *parentModel, msg tea.Msg) *parentModel {
	newModel, _ := p.Update(msg)
	return newModel.(*parentModel)
}

// assertActiveType fails the test if the active model is not of the expected
// type.  Use reflect-free type switches for clarity.
func assertActiveType(t *testing.T, p *parentModel, want string) {
	t.Helper()
	var got string
	switch p.activeModel.(type) {
	case loginModel:
		got = "loginModel"
	case actionModel:
		got = "actionModel"
	case backupModel:
		got = "backupModel"
	case *backupModel:
		got = "backupModel"
	case backupExecModel:
		got = "backupExecModel"
	case listFilesModel:
		got = "listFilesModel"
	case inspectModel:
		got = "inspectModel"
	case restoreModel:
		got = "restoreModel"
	case restoreExecModel:
		got = "restoreExecModel"
	case saveConnectionModel:
		got = "saveConnectionModel"
	default:
		got = "unknown"
	}
	if got != want {
		t.Errorf("activeModel: got %s, want %s", got, want)
	}
}

// ── Routing tests ─────────────────────────────────────────────────────────────

func TestParentModel_LoginFailedMsg_ShowsLoginWithError(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, loginFailedMsg{err: &mockError{"bad password"}})

	assertActiveType(t, p, "loginModel")

	lm := p.activeModel.(loginModel)
	if lm.err == nil || lm.err.Error() != "bad password" {
		t.Errorf("loginModel.err: got %v, want 'bad password'", lm.err)
	}
}

func TestParentModel_LoginDoneMsg_TransitionsToActionModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, loginDoneMsg{host: "localhost:1433", user: "sa", password: "pw"})

	assertActiveType(t, p, "actionModel")
	if p.serverConfig.host != "localhost:1433" {
		t.Errorf("serverConfig.host: got %q, want %q", p.serverConfig.host, "localhost:1433")
	}
}

func TestParentModel_PromptSaveConnectionMsg_TransitionsToSaveModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, promptSaveConnectionMsg{host: "myhost:1433", user: "admin", password: "secret"})

	assertActiveType(t, p, "saveConnectionModel")
	// The server config should be pre-loaded.
	if p.serverConfig.host != "myhost:1433" {
		t.Errorf("serverConfig.host: got %q, want %q", p.serverConfig.host, "myhost:1433")
	}
}

func TestParentModel_ConnectionSavedMsg_TransitionsToActionModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, connectionSavedMsg{host: "localhost:1433", user: "sa", password: "pw"})

	assertActiveType(t, p, "actionModel")
}

func TestParentModel_ConnectionDeletedMsg_TransitionsToLoginModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, connectionDeletedMsg{})

	assertActiveType(t, p, "loginModel")
}

func TestParentModel_GoToActionMsg_TransitionsToActionModel(t *testing.T) {
	p := newTestParent()
	// Start from somewhere other than login to verify we can go back.
	p.activeModel = backupExecModel{err: nil}
	p = sendMsg(p, goToActionMsg{})

	assertActiveType(t, p, "actionModel")
}

func TestParentModel_ActionSelectedMsg_BackupTransitionsToBackupModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, loginDoneMsg{host: "localhost:1433", user: "sa", password: "pw"})
	p = sendMsg(p, actionSelectedMsg{action: "backup"})

	assertActiveType(t, p, "backupModel")
}

func TestParentModel_ActionSelectedMsg_RestoreTransitionsToListFilesModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, loginDoneMsg{host: "localhost:1433", user: "sa", password: "pw"})
	p = sendMsg(p, actionSelectedMsg{action: "restore"})

	assertActiveType(t, p, "listFilesModel")
}

func TestParentModel_DbSelectedMsg_TransitionsToBackupExecModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, dbSelectedMsg{db: "mydb", filename: "mydb.bak"})

	assertActiveType(t, p, "backupExecModel")
}

func TestParentModel_BakFileSelectedMsg_TransitionsToInspectModel(t *testing.T) {
	p := newTestParent()
	p = sendMsg(p, bakFileSelectedMsg{filename: "mydb.bak"})

	assertActiveType(t, p, "inspectModel")
}

func TestParentModel_RestoreBackupMsg_TransitionsToRestoreModel(t *testing.T) {
	p := newTestParent()
	entry := domain.BackupEntry{
		Filename: "mydb.bak",
		MdfFile:  domain.MdfEntry{Name: "mydb"},
		LdfFile:  domain.LdfEntry{Name: "mydb_log"},
	}
	p = sendMsg(p, restoreBackupMsg{fileInfo: entry})

	assertActiveType(t, p, "restoreModel")
}

func TestParentModel_RestoreExecMsg_TransitionsToRestoreExecModel(t *testing.T) {
	p := newTestParent()
	entry := domain.BackupEntry{
		Filename: "mydb.bak",
		MdfFile:  domain.MdfEntry{Name: "mydb"},
		LdfFile:  domain.LdfEntry{Name: "mydb_log"},
	}
	p = sendMsg(p, restoreExecMsg{fileInfo: entry, newDBName: "mydb_restored"})

	assertActiveType(t, p, "restoreExecModel")
}

func TestParentModel_CtrlC_ReturnsQuitCmd(t *testing.T) {
	p := newTestParent()
	_, cmd := p.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Fatal("expected a cmd from ctrl+c, got nil")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg from ctrl+c, got %T", cmd())
	}
}

// ── Full navigation sequence test ─────────────────────────────────────────────

// TestParentModel_FullBackupFlow drives the parent through the complete backup
// navigation path and asserts the active model at each step.
func TestParentModel_FullBackupFlow(t *testing.T) {
	p := newTestParent()
	assertActiveType(t, p, "loginModel")

	// Successful login (already saved connection → loginDoneMsg, not prompt).
	p = sendMsg(p, loginDoneMsg{host: "localhost:1433", user: "sa", password: "pw"})
	assertActiveType(t, p, "actionModel")

	// Select backup action.
	p = sendMsg(p, actionSelectedMsg{action: "backup"})
	assertActiveType(t, p, "backupModel")

	// Database + filename selected.
	p = sendMsg(p, dbSelectedMsg{db: "mydb", filename: "mydb.bak"})
	assertActiveType(t, p, "backupExecModel")

	// Return to action screen.
	p = sendMsg(p, goToActionMsg{})
	assertActiveType(t, p, "actionModel")
}

// TestParentModel_FullRestoreFlow drives the parent through the complete restore
// navigation path.
func TestParentModel_FullRestoreFlow(t *testing.T) {
	p := newTestParent()

	p = sendMsg(p, loginDoneMsg{host: "localhost:1433", user: "sa", password: "pw"})
	assertActiveType(t, p, "actionModel")

	p = sendMsg(p, actionSelectedMsg{action: "restore"})
	assertActiveType(t, p, "listFilesModel")

	p = sendMsg(p, bakFileSelectedMsg{filename: "mydb.bak"})
	assertActiveType(t, p, "inspectModel")

	entry := domain.BackupEntry{
		Filename: "mydb.bak",
		MdfFile:  domain.MdfEntry{Name: "mydb"},
		LdfFile:  domain.LdfEntry{Name: "mydb_log"},
	}
	p = sendMsg(p, restoreBackupMsg{fileInfo: entry})
	assertActiveType(t, p, "restoreModel")

	p = sendMsg(p, restoreExecMsg{fileInfo: entry, newDBName: "mydb_restored"})
	assertActiveType(t, p, "restoreExecModel")

	p = sendMsg(p, goToActionMsg{})
	assertActiveType(t, p, "actionModel")
}

// TestParentModel_SaveConnectionFlow tests the save-prompt detour.
func TestParentModel_SaveConnectionFlow(t *testing.T) {
	p := newTestParent()

	// New connection not yet in the saved list.
	p = sendMsg(p, promptSaveConnectionMsg{host: "newhost:1433", user: "sa", password: "pw"})
	assertActiveType(t, p, "saveConnectionModel")

	// User saves (or skips); both produce connectionSavedMsg.
	p = sendMsg(p, connectionSavedMsg{host: "newhost:1433", user: "sa", password: "pw"})
	assertActiveType(t, p, "actionModel")
}

// TestParentModel_RenderStatusBar checks that RenderStatusBar returns the
// correct breadcrumb for each flow.
func TestParentModel_RenderStatusBar_BackupFlow(t *testing.T) {
	p := newTestParent()
	p.serverConfig.host = "localhost:1433"
	p.activeModel = backupModel{}

	bar := RenderStatusBar(p.serverConfig.host, p.activeModel)
	if bar == "" {
		t.Error("status bar should not be empty for backupModel")
	}
}

func TestParentModel_RenderStatusBar_RestoreFlow(t *testing.T) {
	p := newTestParent()
	p.serverConfig.host = "localhost:1433"
	p.activeModel = listFilesModel{}

	bar := RenderStatusBar(p.serverConfig.host, p.activeModel)
	if bar == "" {
		t.Error("status bar should not be empty for listFilesModel")
	}
}
