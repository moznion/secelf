package internal

import (
	"testing"
)

func TestGenerateContentsKey(t *testing.T) {
	cek, err := generateContentsKey()
	if err != nil {
		t.Fatalf("unexpected err has got [err=%s]", err)
	}
	if len(cek) != 32 {
		t.Errorf("unexpected cek has come [got=%s]", cek)
	}
}
