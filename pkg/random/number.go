package random

import (
	"math/rand"
	"strconv"
)

func RandomPort() string {
	return strconv.Itoa(1024 + rand.Intn(65535-1024))
}
