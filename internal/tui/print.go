package tui

import (

	"strings"
	"fmt"
	"sync"
)

type Result struct {
	Prefix string
	URL    string
	Size   int
	Status int
	Extra string
}
var mu sync.Mutex
func Print(result Result, quiet bool) {
	if quiet {
		fmt.Printf("%d %s %d\n", result.Status, result.URL, result.Size)

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
				Other.Printf("[%s] %s ---> %d %s | %d\n", result.Prefix, result.URL, result.Status, result.Extra,result.Size)
				mu.Unlock()
			} else {
				mu.Lock()
				Other.Printf("[%s] %s ---> %d | %d\n", result.Prefix, result.URL, result.Status,result.Size)
				mu.Unlock()
			}
		}

	}
}
