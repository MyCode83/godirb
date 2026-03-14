package wordlist
import (
	"bufio"
	"strings"
	"fmt"
	"os"
	"errors"
	"godirb/pkg/suggest"
)
type Wordlist struct {
	Wordlist string
	DefaultWordlists []string
}
var ListSlice []string
func (wd *Wordlist) loadReader(list  *os.File) {
	scanner := bufio.NewScanner(list)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		ListSlice = append(ListSlice, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "[X] %v\n", err)
		os.Exit(1)
	}
}
func (WordlistStruct *Wordlist) LoadWordlist() {
	if WordlistStruct.Wordlist == "" {
		fmt.Fprintf(os.Stderr, "[X] Empty Wordlist. Run --help for usage.\n")
		os.Exit(2)
	}
	var list *os.File
	var err error
	if WordlistStruct.Wordlist == "-" {
		list = os.Stdin
		WordlistStruct.loadReader(list)
		return
	}

	list, err = os.Open(WordlistStruct.Wordlist)
	if err != nil {
		if errors.Is(err, os.ErrPermission){
			fmt.Fprintf(os.Stderr, "[X] Error: You don't have permission to read '%s'. Run it as root\n", WordlistStruct.Wordlist)
			os.Exit(2)
		} 
	}
	if err == nil {
		defer list.Close()
		WordlistStruct.loadReader(list)
		return
	}
	WordlistStruct.Wordlist = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(WordlistStruct.Wordlist)), ".txt")
	embedWordlists := []string{"small", "common", "medium", "big", "ports","payloads", "xss"}
	switch  WordlistStruct.Wordlist{
		case "small":
			ListSlice = Small()
			return
		case "common":
			ListSlice = Common()
			return
		case "medium":
			ListSlice = Medium()
		case "big":
			ListSlice = Big()
			return
		case "ports",  "port":
			ListSlice = Ports()
		case "payloads", "payload":
			ListSlice = Payloads()
		case "xss",  "xss-payloads":
			ListSlice = Xss()
		case "lfi":
			ListSlice = Lfi()
		default:
			suggest := suggest.SuggestClosest(2, WordlistStruct.Wordlist, embedWordlists...)
			fmt.Fprintf(os.Stderr,"[X] Wordlist not found\n")
			if  suggest != "" {
				fmt.Printf("Did you mean '%s'?\n", suggest)
			}
			os.Exit(2)
	}
}
