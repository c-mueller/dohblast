package qnamegen

import (
	"bytes"
	"math/rand"
)

const Charset = "abcdefghijklmnopqrstuvwxyz"

var charsetRunes []rune

func init() {
	charsetRunes = bytes.Runes([]byte(Charset))
}

func GenerateRandomString(l int) string {
	id := ""
	for i := 0; i < l; i++ {
		idx := rand.Intn(len(charsetRunes))
		id += string(charsetRunes[idx])
	}
	return id
}
