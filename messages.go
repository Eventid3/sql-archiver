package main

type ConnectionEstablishedMsg struct {
	Config ServerConfig
}

type FileSelectedMsg struct {
	Path string
}

type DatabasesLoadedMsg struct {
	Databases []string
}

// type RestoreOptionsSelectedMsg struct {
// 	Options RestoreOptions
// }

type RestoreCompleteMsg struct {
	Success bool
	Error   error
}

type RestoreProgressMsg struct {
	Percent int
	Message string
}
