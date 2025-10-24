package main

type nextStepMsg struct{}

type formDoneMsg struct {
	user     string
	password string
}

type actionSelectedMsg struct {
	action string
}

type goToActionMsg struct{}
