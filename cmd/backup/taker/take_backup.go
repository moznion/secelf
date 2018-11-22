package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/moznion/secelf/internal"
	"github.com/moznion/secelf/internal/drive"
)

func main() {
	var credentialJSON string
	var tokenJSON string
	var keyArg string
	var sqliteDBPath string
	var backupFilename string

	flag.StringVar(&credentialJSON, "credential-json", "", "[mandatory] credential of Google Drive as JSON string")
	flag.StringVar(&tokenJSON, "token-json", "", "[mandatory] token for accessing to Google Drive as JSON string")
	flag.StringVar(&keyArg, "key", "", "[mandatory] AES key for file encryption (must be 128bit, 192bit or 256bit)")
	flag.StringVar(&sqliteDBPath, "sqlite-db-path", "", "[mandatory] path to SQLite DB file")
	flag.StringVar(&backupFilename, "backup-filename", "", "[mandatory] filename for backup")
	flag.Parse()

	if credentialJSON == "" || tokenJSON == "" || keyArg == "" || sqliteDBPath == "" || backupFilename == "" {
		fmt.Printf("[ERROR] mandatory parameter(s) is/are missing or invalid\n")
		flag.Usage()
		os.Exit(1)
	}

	driveService, err := drive.NewService([]byte(credentialJSON), []byte(tokenJSON))
	if err != nil {
		log.Fatalf("%s", err)
	}

	encrypter, err := internal.NewEncrypter([]byte(keyArg))
	if err != nil {
		log.Fatalf("%s", err)
	}

	contents, err := ioutil.ReadFile(sqliteDBPath)
	if err != nil {
		log.Fatalf("%s", err)
	}

	encrypted, err := encrypter.EncryptBin(contents)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := driveService.UpdateOrPut("", backupFilename, encrypted); err != nil {
		log.Fatalf("%s", err)
	}
}
