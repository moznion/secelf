package internal

import (
	"math/rand"
)

var saltLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var saltLettersLen = len(saltLetters)

func generateSalt(length int) string {
	buff := make([]rune, length)
	for i := range buff {
		buff[i] = saltLetters[rand.Intn(saltLettersLen)]
	}
	return string(buff)
}

func mixKeyAndSalt(masterKey, salt string) []byte {
	saltedKey := make([]byte, len(masterKey))
	for i := range masterKey {
		saltedKey[i] = masterKey[i] ^ salt[i]
	}
	return saltedKey
}
