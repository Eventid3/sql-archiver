package interactive

import "github.com/charmbracelet/lipgloss"

var OuterStyle = lipgloss.NewStyle().Margin(0, 2)

var BorderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	Padding(1, 1).
	BorderForeground(lipgloss.Color("240"))

var DocStyle = lipgloss.NewStyle().Margin(1, 2)

var HeadingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("99")).
	Bold(true).Align(lipgloss.Center)

var ColHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ffffff")).
	Background(lipgloss.Color("99"))

var PurpleBgStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ffffff")).
	Background(lipgloss.Color("#5E5ED2"))

var TableTitleStyle = PurpleBgStyle.
	Padding(0, 1).Margin(0, 2)

var ErrorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
