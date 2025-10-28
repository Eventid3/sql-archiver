package interactive

type nextStepMsg struct{}

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
