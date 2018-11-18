package internal

import (
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

func (r *Registrar) Register(rootDir, fileName string, bin []byte) error {
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

	id, err := r.fileRepo.Put(fileName, encryptedCek)
	if err != nil {
		return err
	}

	extension := filepath.Ext(fileName)
	masqueradeFileName := fmt.Sprintf("%d%s", id, extension)

	enc, err := NewEncrypter(cek)
	if err != nil {
		return err
	}

	encrypted, err := enc.EncryptBin(bin)
	if err != nil {
		return err
	}

	return r.driveService.Put(rootDir, masqueradeFileName, encrypted)
}
