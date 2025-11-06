package interactive

import "github.com/charmbracelet/lipgloss"

var borderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	Margin(1, 2).
	Padding(1, 1).
	BorderForeground(lipgloss.Color("240"))

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var headingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("99")).
	Bold(true).Align(lipgloss.Center)

var colHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ffffff")).
	Background(lipgloss.Color("99"))

var errorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

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

var logo string = `
╔═╗ ┌─┐ ┬     ╔═╗┌─┐┌──┐ ││┐ ┌┌─┌─┐
╚═╗ │ │ │  ── ║═║├┬┘│  ├─┤││ │├─├┬┘
╚═╝ └─┴ ┴─┘   ╝ ╚┘└─└──│ └┴└─┘└─┘└
SQL  ARCHIVER
`
