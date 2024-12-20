package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var currencies = [3]string{EUR, USD, GBP}

// Create a private instance of rand.Rand with a unique source
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomInt generates a random integer between min and max (inclusive).
func RandomInt(min, max int64) int64 {
	if min > max {
		panic("min cannot be greater than max")
	}
	return min + rng.Int63n(max-min+1)
}

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for range n {
		c := alphabet[rng.Intn(k)]
		sb.WriteByte(c)

	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	randCurrenciesIndex := rng.Intn(len(currencies))
	return currencies[randCurrenciesIndex]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
