package tui

import "github.com/charmbracelet/lipgloss"

var (
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	urlStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	redirectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14"))

	clientErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("11"))

	serverErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9"))
)

func statusStyle(status int) lipgloss.Style {
	switch {
	case status >= 200 && status < 300:
		return successStyle
	case status >= 300 && status < 400:
		return redirectStyle
	case status >= 400 && status < 500:
		return clientErrorStyle
	case status >= 500:
		return serverErrorStyle
	default:
		return infoStyle
	}
}
