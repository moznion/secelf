package internal

import (
	"fmt"
	"path/filepath"

	"github.com/moznion/secelf/internal/drive"
	"github.com/moznion/secelf/internal/repository"
)

type Retriever struct {
	key          string
	driveService *drive.Service
	fileRepo     *repository.FileRepository
}

func NewRetriever(key string, fileRepo *repository.FileRepository, driveService *drive.Service) *Retriever {
	return &Retriever{
		key:          key,
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

	enc, err := NewEncrypter(mixKeyAndSalt(r.key, file.Salt))
	if err != nil {
		return nil, err
	}
	return enc.DecryptBin(encrypted), nil
}
