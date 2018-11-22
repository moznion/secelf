package drive

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/moznion/secelf/internal/exception"
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

func (s *Service) Put(rootDir, filename string, content []byte) error {
	file := &gdrive.File{
		Name:    filename,
		Parents: []string{rootDir},
	}

	_, err := s.driveClient.Files.Create(file).Media(bytes.NewReader(content)).Do()
	return err
}

// For only backup
func (s *Service) UpdateOrPut(rootDir, filename string, content []byte) error {
	fileID, err := s.getFileID(rootDir, filename)
	if err != nil && !exception.IsGoogleDriveNotFound(err) {
		return err
	}

	if fileID != "" {
		// update
		file := &gdrive.File{
			Name: filename,
		}
		_, err = s.driveClient.Files.Update(fileID, file).Media(bytes.NewReader(content)).Do()
		return err
	}

	// create new
	file := &gdrive.File{
		Name: filename,
	}
	if rootDir != "" {
		file.Parents = []string{rootDir}
	}

	_, err = s.driveClient.Files.Create(file).Media(bytes.NewReader(content)).Do()
	return err
}

func (s *Service) Get(rootDir, filename string) ([]byte, error) {
	fileID, err := s.getFileID(rootDir, filename)
	if err != nil {
		return nil, err
	}

	resp, err := s.driveClient.Files.Get(fileID).Download()
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

func (s *Service) getFileID(rootDir, filename string) (string, error) {
	req := s.driveClient.Files.List().Fields("nextPageToken, files(id)").Q(fmt.Sprintf("name = '%s'", filename))
	if rootDir != "" {
		req = req.Q(fmt.Sprintf("'%s' in parents", rootDir))
	}

	fileList, err := req.Do()
	if err != nil {
		return "", err
	}

	files := fileList.Files
	if len(files) <= 0 {
		return "", exception.BuildGoogleDriveNotFound(fmt.Sprintf("not found [rootDir=%s, filename=%s]", rootDir, filename))
	}

	return files[0].Id, nil
}
