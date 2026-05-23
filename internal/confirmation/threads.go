package confirmation

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ThreadsConfirmation(msg string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg + " (Y/N): ")
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes" || input == "si"

}
