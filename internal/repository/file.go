package repository

import (
	"database/sql"
	"fmt"

	// driver for SQLite3
	_ "github.com/mattn/go-sqlite3"
)

const fileTableName = "file"

type FileRepository struct {
	dbPath string
}

func NewFileRepository(dbPath string) *FileRepository {
	return &FileRepository{
		dbPath: dbPath,
	}
}

func (repo *FileRepository) Put(filename string, encryptedCek []byte) (int64, error) {
	db, err := sql.Open("sqlite3", repo.dbPath)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (filename, encrypted_cek) VALUES (?, ?)", fileTableName))
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(filename, encryptedCek)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

type FileRow struct {
	ID           int64  `json:"id"`
	Filename     string `json:"filename"`
	EncryptedCek []byte `json:"encrypted_cek"`
}

func (repo *FileRepository) Single(id int64) (*FileRow, error) {
	db, err := sql.Open("sqlite3", repo.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("SELECT filename, encrypted_cek FROM %s WHERE id = ?", fileTableName))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var filename string
	var encryptedCek []byte
	err = stmt.QueryRow(id).Scan(&filename, &encryptedCek)
	if err != nil {
		return nil, err
	}

	return &FileRow{
		ID:           id,
		Filename:     filename,
		EncryptedCek: encryptedCek,
	}, nil
}

func (repo *FileRepository) Search(q string) ([]*FileRow, error) {
	db, err := sql.Open("sqlite3", repo.dbPath)
	if err != nil {
		return []*FileRow{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, filename FROM %s WHERE filename LIKE ? ORDER BY id DESC", fileTableName))
	if err != nil {
		return []*FileRow{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + q + "%")
	if err != nil {
		return []*FileRow{}, err
	}
	defer rows.Close()

	results := make([]*FileRow, 0)
	for rows.Next() {
		var id int64
		var filename string
		err = rows.Scan(&id, &filename)
		if err != nil {
			continue
		}
		results = append(results, &FileRow{
			ID:       id,
			Filename: filename,
		})
	}

	return results, nil
}
