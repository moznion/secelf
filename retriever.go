package secelf

import (
	"fmt"
	"path/filepath"

	"github.com/moznion/secelf/internal/drive"
	"github.com/moznion/secelf/repository"
)

type Retriever struct {
	enc          *Encrypter
	driveService *drive.Service
	fileRepo     *repository.FileRepository
}

func NewRetriever(enc *Encrypter, fileRepo *repository.FileRepository, driveService *drive.Service) *Retriever {
	return &Retriever{
		enc:          enc,
		driveService: driveService,
		fileRepo:     fileRepo,
	}
}

func (r *Retriever) Retrieve(id int64, rootDir string) ([]byte, error) {
	file, err := r.fileRepo.Single(id)
	if err != nil {
		return nil, err
	}

	// TODO
	actualFileName := file.FileName
	masqueradeFileName := fmt.Sprintf("%d%s", id, filepath.Ext(actualFileName))

	encrypted, err := r.driveService.Get(rootDir, masqueradeFileName)
	if err != nil {
		return nil, err
	}

	return r.enc.DecryptBin(encrypted), nil
}
