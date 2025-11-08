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

type BakFileInfo struct {
	filename, mdfName, ldfName, mdfSize, mdfBackupSize, ldfSize string
}

type restoreBackupMsg struct {
	fileInfo BakFileInfo
}

type restoreExecMsg struct {
	fileInfo  BakFileInfo
	newDBName string
}
