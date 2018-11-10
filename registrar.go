package secelf

import (
	"fmt"
	"path/filepath"

	"github.com/moznion/secelf/internal/drive"
	"github.com/moznion/secelf/repository"
)

type Registrar struct {
	enc          *Encrypter
	driveService *drive.Service
	fileRepo     *repository.FileRepository
}

func NewRegistrar(enc *Encrypter, fileRepo *repository.FileRepository, driveService *drive.Service) *Registrar {
	return &Registrar{
		enc:          enc,
		driveService: driveService,
		fileRepo:     fileRepo,
	}
}

func (r *Registrar) Register(rootDir, fileName string, bin []byte) error {
	id, err := r.fileRepo.Put(fileName)
	if err != nil {
		return err
	}

	extension := filepath.Ext(fileName)
	masqueradeFileName := fmt.Sprintf("%d%s", id, extension)

	encrypted, err := r.enc.EncryptBin(bin)
	if err != nil {
		return err
	}

	return r.driveService.Put(rootDir, masqueradeFileName, encrypted)
}
