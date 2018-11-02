package secelf

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Encrypter struct {
	block cipher.Block
}

func NewEncrypter(key []byte) (*Encrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Encrypter{
		block: block,
	}, nil
}

func (enc *Encrypter) EncryptBin(raw []byte) ([]byte, error) {
	encrypted := make([]byte, aes.BlockSize+len(raw))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	encryptStream := cipher.NewCTR(enc.block, iv)
	encryptStream.XORKeyStream(encrypted[aes.BlockSize:], raw)
	return encrypted, nil
}

func (enc *Encrypter) DecryptBin(encrypted []byte) []byte {
	decryptStream := cipher.NewCTR(enc.block, encrypted[:aes.BlockSize])

	decrypted := make([]byte, len(encrypted[aes.BlockSize:]))
	decryptStream.XORKeyStream(decrypted, encrypted[aes.BlockSize:])
	return decrypted
}
