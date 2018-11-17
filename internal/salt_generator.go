package internal

import (
	"crypto/rand"
	"math/big"
)

var saltLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var saltLettersLen = len(saltLetters)

func generateSalt(length int) (string, error) {
	buff := make([]rune, length)
	for i := range buff {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(saltLettersLen)))
		if err != nil {
			return "", err
		}
		buff[i] = saltLetters[n.Sign()]
	}
	return string(buff), nil
}

func mixKeyAndSalt(masterKey, salt string) []byte {
	saltedKey := make([]byte, len(masterKey))
	for i := range masterKey {
		saltedKey[i] = masterKey[i] ^ salt[i]
	}
	return saltedKey
}
