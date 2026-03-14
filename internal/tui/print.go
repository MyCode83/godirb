package tui

import (

	"strings"
	"fmt"
)

type Result struct {
	Prefix string
	URL    string
	Size   int
	Status int
	Extra string
}
func Print(result Result, quiet bool) {
	if quiet {
		fmt.Printf("%d %s %d\n", result.Status, result.URL, result.Size)

	} else {

		isFile := strings.Contains(result.Prefix, "FILE")
		switch isFile {
		case true:
			File.Printf("[%s] %s ---> %d | %d\n", result.Prefix, result.URL, result.Status, result.Size)
		default:
			if strings.TrimSpace(result.Extra) != "" {
				Other.Printf("[%s] %s ---> %d %s | %d\n", result.Prefix, result.URL, result.Status, result.Extra,result.Size)
			} else {
				Other.Printf("[%s] %s ---> %d | %d\n", result.Prefix, result.URL, result.Status,result.Size)
			}
		}

	}
}
