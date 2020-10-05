package shared

import (
	"math/rand"
	"time"
)

func GetRandomString(charset string, l int64) string {
	b := make([]byte, l)
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
