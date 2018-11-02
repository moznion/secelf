package secelf

import (
	"testing"
)

func TestEncDecRoundTripShouldSuccessfully(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // 32 bytes
	enc, err := NewEncrypter(key)
	if err != nil {
		t.Fatalf("got unexpected error when creating a new encrypter")
	}

	raw := []byte("this is a pen")

	encrypted, err := enc.EncryptBin(raw)
	if err != nil {
		t.Fatalf("got unexpected error when encryption")
	}

	decrypted := enc.DecryptBin(encrypted)
	if string(decrypted) != string(raw) {
		t.Fatalf("decrypted result is wrong; something broken?")
	}
}
