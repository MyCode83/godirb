package assemble

import (
	"fmt"

	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/core"
)

const banner string = (`		                     
   ____ _  ____   ____/ /   (_)   _____   / /_
  / __  / / __ \ / __  /   / /   / ___/  / __ \
 / /_/ / / /_/ // /_/ /   / /   / /     / /_/ /
 \__  /  \____/ \____/   /_/   /_/     /_____/
/____/
`)

func PrintBanner(cfg cli.Config, mode core.Mode, quietOutput bool) {
	if !cfg.Quiet && !(quietOutput && cfg.Output == "") {
		fmt.Printf(banner)
		fmt.Println("\n------------------")
		fmt.Println("[*] Url: ", cfg.BaseURL)
		fmt.Println("[*] Method: ", cfg.Method)
		fmt.Println("[*] Threads: ", cfg.Threads)
		fmt.Println("[*] Timeout: ", cfg.Timeout)
		fmt.Println("[*] Delay: ", cfg.Delay)
		fmt.Println("[*] UAs: ", len(cfg.UserAgent))
		fmt.Print("[*] Mode: ")
		switch mode {
		case core.ModeDir:
			fmt.Print("Dir\n")
		case core.ModeFuzz:
			fmt.Print("Fuzz\n")
		case core.ModePort:
			fmt.Print("Port\n")
		}
		fmt.Printf("------------------\n\n")
	}
}
