package wordlist

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/MyCode83/godirb/pkg/suggest"
	"io"
	"os"
	"strings"
)

type Wordlist struct {
	Wordlist         string
	DefaultWordlists []string
}

var listSlice []string

func (wd *Wordlist) loadReader(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "\r", "")
		line = strings.ReplaceAll(scanner.Text(), "\t", "")

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fmt.Println(line)
		listSlice = append(listSlice, line)
	}
}
func (WordlistStruct *Wordlist) LoadWordlist() []string {
	if WordlistStruct.Wordlist == "" {
		fmt.Fprintf(os.Stderr, "[X] Empty Wordlist. Run --help for usage.\n")
		os.Exit(2)
	}
	var list *os.File
	var err error
	fmt.Println()
	if WordlistStruct.Wordlist == "-" {
		if WordlistStruct.Wordlist == "-" {

			WordlistStruct.loadReader(os.Stdin)
			return listSlice
		}
	}

	list, err = os.Open(WordlistStruct.Wordlist)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			fmt.Fprintf(os.Stderr, "[X] Error: You don't have permission to read '%s'. Run it as root\n", WordlistStruct.Wordlist)
			os.Exit(2)
		}
	}
	if err == nil {
		defer list.Close()
		WordlistStruct.loadReader(list)
		return listSlice
	}
	WordlistStruct.Wordlist = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(WordlistStruct.Wordlist)), ".txt")
	embedWordlists := []string{"small", "common", "medium", "big", "ports", "payloads", "xss", "lfi"}
	switch WordlistStruct.Wordlist {
	case "small":
		listSlice = Small()
	case "common":
		listSlice = Common()
	case "medium":
		listSlice = Medium()
	case "big":
		listSlice = Big()
	case "ports", "port":
		listSlice = Ports()
	case "payloads", "payload":
		listSlice = Payloads()
	case "xss", "xss-payloads":
		listSlice = Xss()
	case "lfi":
		listSlice = Lfi()
	default:
		suggest := suggest.SuggestClosest(2, WordlistStruct.Wordlist, embedWordlists...)
		fmt.Fprintf(os.Stderr, "[X] Wordlist not found\n")
		if suggest != "" {
			fmt.Fprintf(os.Stderr, "Did you mean '%s'?\n", suggest)
		}
		os.Exit(2)
	}
	return listSlice

}
