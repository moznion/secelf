package drive

import (
	"bytes"
	"fmt"
	"io/ioutil"

	gdrive "google.golang.org/api/drive/v3"
)

type Service struct {
	driveClient *gdrive.Service
}

func NewService(credentialJSON []byte, tokenJSON []byte) (*Service, error) {
	driveClient, err := makeDriveClient(credentialJSON, tokenJSON)
	if err != nil {
		return nil, err
	}

	return &Service{
		driveClient: driveClient,
	}, nil
}

func (s *Service) Put(rootDir, fileName string, content []byte) error {
	file := &gdrive.File{
		Name:    fileName,
		Parents: []string{rootDir},
	}

	_, err := s.driveClient.Files.Create(file).Media(bytes.NewReader(content)).Do()
	return err
}

func (s *Service) Get(rootDir, fileName string) ([]byte, error) {
	fileList, err := s.driveClient.Files.List().Fields("nextPageToken, files(id)").Q(fmt.Sprintf("'%s' in parents", rootDir)).Do()
	if err != nil {
		return nil, err
	}

	files := fileList.Files
	if len(files) <= 0 {
		return nil, fmt.Errorf("not found [rootDir=%s, fileName=%s]", rootDir, fileName)
	}

	resp, err := s.driveClient.Files.Get(files[0].Id).Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bin, nil
}
