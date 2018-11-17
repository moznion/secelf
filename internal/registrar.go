package internal

import (
	"fmt"
	"path/filepath"

	"github.com/moznion/secelf/internal/drive"
	"github.com/moznion/secelf/internal/repository"
)

type Registrar struct {
	key          string
	driveService *drive.Service
	fileRepo     *repository.FileRepository
}

func NewRegistrar(key string, fileRepo *repository.FileRepository, driveService *drive.Service) *Registrar {
	return &Registrar{
		key:          key,
		driveService: driveService,
		fileRepo:     fileRepo,
	}
}

func (r *Registrar) Register(rootDir, fileName string, bin []byte) error {
	salt, err := generateSalt(len(r.key))
	if err != nil {
		return err
	}

	id, err := r.fileRepo.Put(fileName, salt)
	if err != nil {
		return err
	}

	extension := filepath.Ext(fileName)
	masqueradeFileName := fmt.Sprintf("%d%s", id, extension)

	enc, err := NewEncrypter(mixKeyAndSalt(r.key, salt))
	if err != nil {
		return err
	}

	encrypted, err := enc.EncryptBin(bin)
	if err != nil {
		return err
	}

	return r.driveService.Put(rootDir, masqueradeFileName, encrypted)
}
