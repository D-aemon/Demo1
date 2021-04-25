package util

import (
	"math/rand"
	"time"
)

//随机生成一个用户名
func RandomString(n int) string {
	var letter = []byte("qwertyuiopassdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXVBNM")
	res := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range res {
		res[i] = letter[rand.Intn(len(letter))]
	}
	return string(res)
}
