package internal

import (
	"fmt"
	"path/filepath"

	"github.com/moznion/secelf/internal/drive"
	"github.com/moznion/secelf/internal/repository"
)

type Retriever struct {
	kek          []byte
	driveService *drive.Service
	fileRepo     *repository.FileRepository
}

func NewRetriever(key []byte, fileRepo *repository.FileRepository, driveService *drive.Service) *Retriever {
	return &Retriever{
		kek:          key,
		driveService: driveService,
		fileRepo:     fileRepo,
	}
}

func (r *Retriever) Retrieve(id int64, rootDir string) ([]byte, error) {
	file, err := r.fileRepo.Single(id)
	if err != nil {
		return nil, err
	}

	actualFilename := file.Filename
	masqueradeFilename := fmt.Sprintf("%d%s", id, filepath.Ext(actualFilename))

	encrypted, err := r.driveService.Get(rootDir, masqueradeFilename)
	if err != nil {
		return nil, err
	}

	contentsKeyEncrypter, err := NewEncrypter(r.kek)
	if err != nil {
		return nil, err
	}
	cek := contentsKeyEncrypter.DecryptBin(file.EncryptedCek)

	enc, err := NewEncrypter(cek)
	if err != nil {
		return nil, err
	}
	return enc.DecryptBin(encrypted), nil
}
