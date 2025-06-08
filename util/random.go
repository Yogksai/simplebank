package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Randomint(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}
func RandomBalance() int64 {
	return Randomint(0, 1000000)
}
func RandomCurrency() string {
	currencies := []string{"TNG", "RUB", "SOM"}
	return currencies[rand.Intn(len(currencies))]
}
