package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

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

func ConfigureColor(noColor bool, quiet bool) {
	if colorDisabled(noColor, quiet) {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
}

func colorDisabled(noColor bool, quiet bool) bool {
	return noColor || quiet || envIsSet("NO_COLOR") || envIsSet("GODIRB_NO_COLOR")
}

func envIsSet(name string) bool {
	_, ok := os.LookupEnv(name)
	return ok
}
