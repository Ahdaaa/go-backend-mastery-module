package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// random float between 0.0 and max
func RandomFloat(max float64) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64() * max
}

// this will generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() float64 {
	return RandomFloat(1000.0)
}

func RandomCurrency() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	currencies := []string{"EUR", "USD", "IDR"}
	n := len(currencies)
	return currencies[r.Intn(n)]
}
