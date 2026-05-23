package tui

import (
	"fmt"
	"strings"
	"sync"

	"godirb/internal/core"
	"godirb/internal/output"
)

var mu sync.Mutex

func Print(result core.Result, quiet bool) {
	if quiet {
		fmt.Println(output.FormatTextResult(result, quiet))

	} else {

		isFile := strings.Contains(result.Prefix, "FILE")
		switch isFile {
		case true:
			mu.Lock()
			File.Printf("[%s] %s ---> %d | %d\n", result.Prefix, result.URL, result.Status, result.Size)
			mu.Unlock()
		default:
			if strings.TrimSpace(result.Extra) != "" {
				mu.Lock()
				Other.Println(output.FormatTextResult(result, quiet))
				mu.Unlock()
			} else {
				mu.Lock()
				Other.Println(output.FormatTextResult(result, quiet))
				mu.Unlock()
			}
		}

	}
}
