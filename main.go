package secelf

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/moznion/secelf/internal"
	"github.com/moznion/secelf/internal/drive"
	"github.com/moznion/secelf/internal/repository"
)

func Run(args []string) {
	var port int64
	var credentialJSON string
	var tokenJSON string
	var key string
	var rootDirID string
	var sqliteDBPath string

	flag.Int64Var(&port, "port", -1, "[mandatory] port for listen")
	flag.StringVar(&credentialJSON, "credential-json", "", "[mandatory] credential of Google Drive as JSON string")
	flag.StringVar(&tokenJSON, "token-json", "", "[mandatory] token for accessing to Google Drive as JSON string")
	flag.StringVar(&key, "key", "", "[mandatory] AES key for file encryption (must be 128bit, 192bit or 256bit)")
	flag.StringVar(&rootDirID, "root-dir-id", "", "[mandatory] identifier fo root directory for storing files")
	flag.StringVar(&sqliteDBPath, "sqlite-db-path", "", "[mandatory] path to SQLite DB file")
	flag.Parse()

	if port < 0 || credentialJSON == "" || tokenJSON == "" || key == "" || rootDirID == "" || sqliteDBPath == "" {
		fmt.Printf("[ERROR] mandatory parameter(s) is/are missing or invalid\n")
		flag.Usage()
		os.Exit(1)
	}

	driveService, err := drive.NewService([]byte(credentialJSON), []byte(tokenJSON))
	if err != nil {
		log.Fatalf("%s", err)
	}

	fileRepo := repository.NewFileRepository(sqliteDBPath)

	register := internal.NewRegistrar(key, fileRepo, driveService)
	retriever := internal.NewRetriever(key, fileRepo, driveService)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			w.Write([]byte("method not allowed"))
			return
		}

		tmpl := template.Must(template.ParseFiles("./webui/index.html"))
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			log.Printf("[ERROR] %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}
	})

	r.HandleFunc("/files/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
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

		content, err := retriever.Retrieve(id, rootDirID)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		w.WriteHeader(200)
		w.Write(content)
	})

	r.HandleFunc("/webui/dist/{file}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	r.HandleFunc("/api/files", func(w http.ResponseWriter, r *http.Request) {
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

		err = register.Register(rootDirID, fileName, bin)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		w.WriteHeader(201)
		w.Write([]byte(fmt.Sprintf("created: %s", fileName)))
	})

	r.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			w.Write([]byte("method not allowed"))
			return
		}

		q := r.URL.Query().Get("q")
		records, err := fileRepo.Search(q)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		result, _ := json.Marshal(records)

		w.WriteHeader(200)
		w.Write(result)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("server start: %s", addr)
	log.Fatal(srv.ListenAndServe())
}
