package interactive

import "github.com/Eventid3/sql-archiver/domain"

type nextStepMsg struct{}

type loginFailedMsg struct {
	err error
}

type loginDoneMsg struct {
	host     string
	user     string
	password string
}

// promptSaveConnectionMsg is emitted after a successful login when the
// connection is not already in the saved list. The user is asked to name and
// save it before proceeding to the action screen.
type promptSaveConnectionMsg struct {
	host     string
	user     string
	password string
}

// connectionSavedMsg is emitted by saveConnectionModel after the connection
// has been written to disk. The embedded loginDoneMsg carries the credentials
// so parentModel can proceed to the action screen.
type connectionSavedMsg struct {
	host     string
	user     string
	password string
}

// connectionDeletedMsg is emitted after a saved connection has been removed
// from disk. parentModel re-renders the login screen with the refreshed list.
type connectionDeletedMsg struct{}

type actionSelectedMsg struct {
	action string
}

type goToActionMsg struct{}

type dbSelectedMsg struct {
	db       string
	filename string
}

type bakFileSelectedMsg struct {
	filename string
}

type restoreBackupMsg struct {
	fileInfo domain.BackupEntry
}

type restoreExecMsg struct {
	fileInfo  domain.BackupEntry
	newDBName string
}
