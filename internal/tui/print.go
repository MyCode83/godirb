package tui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/output"
)

var mu sync.Mutex

func Print(result core.Result, quiet bool) {
	var line string

	if quiet {
		line = output.FormatTextResult(result, true)
	} else {
		line = renderResult(result)
	}

	mu.Lock()
	fmt.Println(line)
	mu.Unlock()
}

func renderResult(result core.Result) string {
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
