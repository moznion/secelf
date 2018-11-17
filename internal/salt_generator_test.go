package internal

import (
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	l := 32
	salt, err := generateSalt(l)
	if err != nil {
		t.Fatalf("unexpected err has got [err=%s]", err)
	}
	if len(salt) != l {
		t.Errorf("unexpected salt has come [got=%s]", salt)
	}
}

func TestMixKeyAndSalt(t *testing.T) {
	saltedKey := mixKeyAndSalt("bob", "sam")
	if len(saltedKey) != 3 {
		t.Fatalf("unexpected saltedKey length has come [got=%v]", saltedKey)
	}
	if saltedKey[0] != 17 || saltedKey[1] != 14 || saltedKey[2] != 15 {
		t.Fatalf("unexpected saltedKey has come [got=%v]", saltedKey)
	}
}
