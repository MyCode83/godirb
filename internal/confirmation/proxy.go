package confirmation

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ProxyConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(" (y/N): ")
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes" || input == "si"
}
