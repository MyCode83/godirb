package random
import "math/rand"
func RandChoice(slice []string) string{
	return slice[rand.Intn(len(slice))]
}