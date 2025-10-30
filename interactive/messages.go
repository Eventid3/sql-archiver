package interactive

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
	filename, mdfName, ldfName string
}

type restoreExecMsg struct {
	filename, mdfName, ldfName, newDBName string
}
