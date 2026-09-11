package ui

import (
	"fmt"
	"strings"

	"github.com/MyCode83/godirb/internal/core"
)

func RenderResult(result core.Result) string {
	kind := infoStyle.Render(
		fmt.Sprintf("%-7s", strings.ToUpper(result.Kind)),
	)

	status := statusStyle(result.Status).Render(
		fmt.Sprintf("%3d", result.Status),
	)

	size := infoStyle.Render(
		fmt.Sprintf("%8d B", result.Size),
	)

	url := urlStyle.Render(result.URL)

	line := fmt.Sprintf("%s  %s  %s  %s", kind, status, size, url)

	if extra := strings.TrimSpace(result.Error); extra != "" {
		line += "  " + infoStyle.Render(extra)
	}

	return line
}
