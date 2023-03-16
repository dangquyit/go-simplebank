package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"

func init() {
	rand.Seed(int64(time.Now().UnixNano()))

}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // min ->max
}

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

func RandomUserName() string {
	return RandomOwner() + strconv.Itoa(int(RandomInt(1000, 9999)))
}

func RandomPassword() string {
	return RandomString(11)
}

func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

func RandomAccountNumber() int64 {
	return int64(RandomInt(1000000000, 9999999999))
}

func RandomCurrency() string {
	currencies := []string{"USD", "VND", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	prefix := RandomString(6)
	suffixs := []string{"@gmail.com", "@microsoft.com", "@icloud.com"}
	n := len(suffixs)

	return prefix + suffixs[rand.Intn(n)]
}
