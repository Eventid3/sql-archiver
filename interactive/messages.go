package interactive

import "github.com/Eventid3/sql-archiver/domain"

type nextStepMsg struct{}

type loginFailedMsg struct {
	err error
}

type loginDoneMsg struct {
	container string
	user      string
	password  string
}

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
