package internal

import (
	"crypto/rand"
	"math/big"
)

var max256int, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819967", 10)

//                                        ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ (1 << 255) - 1

func generateContentsKey() ([]byte, error) {
	n, err := rand.Int(rand.Reader, max256int)
	if err != nil {
		return nil, err
	}
	return n.Bytes(), nil
}
