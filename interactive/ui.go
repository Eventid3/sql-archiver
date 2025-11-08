package interactive

import tea "github.com/charmbracelet/bubbletea"

var buildingBlocks = `
╔
╗
╚
╝
═
║
┬
┴
├
┤
┼
┌
┐
└
┘
─
│
`

var Logo string = `
╔═╗ ┌─┐ ┬     ╔═╗┌─┐┌──┐ ││┐ ┌┌─┌─┐
╚═╗ │ │ │  ── ║═║├┬┘│  ├─┤││ │├─├┬┘
╚═╝ └─┴ ┴─┘   ╝ ╚┘└─└──│ └┴└─┘└─┘└
SQL  ARCHIVER
`

func RenderStatusBar(container string, model tea.Model) string {
	currentFlow := ""
	switch model.(type) {
	case backupModel, backupExecModel:
		currentFlow = " > Backup"
	case restoreModel, restoreExecModel, listFilesModel, inspectModel:
		currentFlow = " > Restore"
	}
	return ColHeaderStyle.Padding(0, 1).Bold(true).Render("Container: " + container + currentFlow)
}
