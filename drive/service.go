package drive

import (
	"bytes"

	gdrive "google.golang.org/api/drive/v3"
)

type Service struct {
	driveClient *gdrive.Service
}

func NewService(driveClient *gdrive.Service) *Service {
	return &Service{
		driveClient: driveClient,
	}
}

func (s *Service) Put(rootDir, fileName string, content []byte) error {
	file := &gdrive.File{
		Name:    fileName,
		Parents: []string{rootDir},
	}

	_, err := s.driveClient.Files.Create(file).Media(bytes.NewReader(content)).Do()
	return err
}
