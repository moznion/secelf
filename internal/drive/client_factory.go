package drive

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gdrive "google.golang.org/api/drive/v3"
)

func MakeDriveClient(credentialJSON []byte, tokenJSON []byte) (*gdrive.Service, error) {
	config, err := google.ConfigFromJSON(credentialJSON, gdrive.DriveScope)
	if err != nil {
		return nil, err
	}

	httpClient, err := getHTTPClient(config, tokenJSON)
	if err != nil {
		return nil, err
	}

	driveClient, err := gdrive.New(httpClient)
	if err != nil {
		return nil, err
	}
	return driveClient, nil
}

func getHTTPClient(config *oauth2.Config, tokenJSON []byte) (*http.Client, error) {
	token := &oauth2.Token{}
	err := json.Unmarshal(tokenJSON, token)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), token), nil
}
