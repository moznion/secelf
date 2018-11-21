package internal

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/moznion/secelf/internal/drive"
	"github.com/moznion/secelf/internal/repository"
)

type Registrar struct {
	kek          []byte
	driveService *drive.Service
	fileRepo     *repository.FileRepository
}

func NewRegistrar(key []byte, fileRepo *repository.FileRepository, driveService *drive.Service) *Registrar {
	return &Registrar{
		kek:          key,
		driveService: driveService,
		fileRepo:     fileRepo,
	}
}

func (r *Registrar) Register(rootDir, filename string, bin []byte) error {
	cek, err := generateContentsKey()
	if err != nil {
		return err
	}

	contentsKeyEncrypter, err := NewEncrypter(r.kek)
	if err != nil {
		return err
	}
	encryptedCek, err := contentsKeyEncrypter.EncryptBin(cek)
	if err != nil {
		return err
	}

	id, err := r.fileRepo.Put(filename, encryptedCek)
	if err != nil {
		return err
	}

	extension := filepath.Ext(filename)
	masqueradeFilename := fmt.Sprintf("%d%s", id, extension)

	enc, err := NewEncrypter(cek)
	if err != nil {
		return err
	}

	encrypted, err := enc.EncryptBin(bin)
	if err != nil {
		return err
	}

	err = r.driveService.Put(rootDir, masqueradeFilename, encrypted)
	if err != nil {
		return err
	}

	return r.verifyUpload(rootDir, masqueradeFilename, id, encrypted)
}

func (r *Registrar) verifyUpload(rootDir, masqueradeFilename string, id int64, encrypted []byte) error {
	got, err := r.driveService.Get(rootDir, masqueradeFilename)
	if err != nil {
		return err
	}

	if bytes.Compare(got, encrypted) != 0 {
		return fmt.Errorf("failed upload: verifycation is not passed [id=%d]", id)
	}

	return nil
}
