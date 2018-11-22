package exception

import (
	"fmt"
	"strings"
)

const GoogleDriveNotFound = "ERR001"

func BuildGoogleDriveNotFound(errBody string) error {
	return fmt.Errorf("%s: %s", GoogleDriveNotFound, errBody)
}

func IsGoogleDriveNotFound(err error) bool {
	return strings.HasPrefix(err.Error(), GoogleDriveNotFound)
}
