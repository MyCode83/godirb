package wordlist
import (
	"bufio"
	"strings"
	"fmt"
	"os"
	"errors"
	"godirb/pkg/suggest"
	"io"
)
type Wordlist struct {
	Wordlist string
	DefaultWordlists []string
}
var ListSlice []string
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
		ListSlice = append(ListSlice, line)
	}
}
func (WordlistStruct *Wordlist) LoadWordlist() {
	if WordlistStruct.Wordlist == "" {
		fmt.Fprintf(os.Stderr, "[X] Empty Wordlist. Run --help for usage.\n")
		os.Exit(2)
	}
	var list *os.File
	var err error
	fmt.Println()
	if WordlistStruct.Wordlist == "-" {
		if WordlistStruct.Wordlist == "-" {
			fmt.Println("STDIN")
			
			WordlistStruct.loadReader(os.Stdin)
			return
		}
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
	embedWordlists := []string{"small", "common", "medium", "big", "ports","payloads", "xss", "lfi"}
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
				fmt.Fprintf(os.Stderr, "Did you mean '%s'?\n", suggest)
			}
			os.Exit(2)
	}
}
