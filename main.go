package secelf

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/moznion/secelf/drive"
	"github.com/moznion/secelf/repository"
)

func Run(args []string) {
	// TODO these should be extract to parameter?
	credentialJSON := os.Getenv("CREDENTIAL_JSON")
	if credentialJSON == "" {
		log.Fatalf(`ENV[CREDENTIAL_JSON] is mandatory variable, however it is missing`)
	}

	tokenJSON := os.Getenv("TOKEN_JSON")
	if tokenJSON == "" {
		log.Fatalf(`ENV[TOKEN_JSON] is mandatory variable, however it is missing`)
	}

	parentDirID := os.Getenv("PARENT_DIR_ID")
	if parentDirID == "" {
		log.Fatalf(`ENV[PARENT_DIR_ID] is mandatory variable, however it is missing`)
	}

	key := os.Getenv("KEY")
	if key == "" {
		log.Fatalf(`ENV[KEY] is mandatory variable, however it is missing`)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatalf(`ENV[DB_PATH] is mandatory variable, however it is missing`)
	}

	rootDir := os.Getenv("ROOT_DIR")
	if rootDir == "" {
		log.Fatalf(`ENV[ROOT_DIR] is mandatory variable, however it is missing`)
	}

	enc, err := NewEncrypter([]byte(key))
	if err != nil {
		log.Fatalf("%s", err)
	}

	driveClient, err := drive.MakeDriveClient([]byte(credentialJSON), []byte(tokenJSON))
	if err != nil {
		log.Fatalf("%s", err)
	}
	driveService := drive.NewService(driveClient)

	fileRepo := repository.NewFileRepository(dbPath)

	register := NewRegistrar(enc, fileRepo, driveService)
	retriever := NewRetriever(enc, fileRepo, driveService)

	r := mux.NewRouter()

	r.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			w.Write([]byte("method not allowed"))
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("[ERROR] %s", err)
			w.WriteHeader(400)
			w.Write([]byte("invalid request"))
			return
		}

		defer file.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, file)
		bin := buf.Bytes()

		fileName := header.Filename

		err = register.Register(rootDir, fileName, bin)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		w.WriteHeader(201)
		w.Write([]byte("ok"))
	})

	r.HandleFunc("/file/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			w.Write([]byte("method not allowed"))
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}

		content, err := retriever.Retrieve(id, rootDir)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		w.WriteHeader(200)
		w.Write(content)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:29292",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
