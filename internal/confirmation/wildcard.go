package confirmation
import (
	"bufio"
	"os"
	"strings"
)
func WildcardConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes" || input == "si" || input == "Y"

}