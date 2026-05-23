package random

import "math/rand"

func RandomString(nlenght int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	buffer := make([]byte, nlenght)
	for index := range buffer {
		buffer[index] = letters[rand.Intn(len(letters))]
	}
	return string(buffer)
}
