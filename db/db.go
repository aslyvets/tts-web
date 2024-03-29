package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
	"tts-web/model"
	"tts-web/service"
)

func InitializeDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("TTS_DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func SaveTTSAudio(db *sql.DB, record model.TTSRecord) error {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return err
	}

	if isTitleExists(db, record.Title) {
		record.Title = service.IncrementStringEnd(record.Title)
	}

	stmt, err := db.Prepare(
		"INSERT INTO tts_responses(id, title, text_input, audio_data, created_at) VALUES($1, $2, $3, $4, $5)",
	)
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, record.Title, record.InputText, record.AudioContent, time.Now())
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}

func FetchTTSAudioByID(db *sql.DB, id string) ([]byte, error) {
	var audio []byte
	err := db.QueryRow("SELECT audio_data FROM tts_responses WHERE id = $1", id).Scan(&audio)
	if err != nil {
		return nil, err
	}
	return audio, nil
}

func DeleteTTSRecordByID(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM tts_responses WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func FetchAllTTSRecords(db *sql.DB) ([]model.TTSListRecord, error) {
	rows, err := db.Query("SELECT id, title, text_input FROM tts_responses ORDER BY created_at DESC")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.TTSListRecord
	for rows.Next() {
		var record model.TTSListRecord
		if err := rows.Scan(&record.Id, &record.Title, &record.Text); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// isTitleExists checks if a record with the given title exists in the database.
func isTitleExists(db *sql.DB, title string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM tts_responses WHERE title = $1)"

	err := db.QueryRow(query, title).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	return exists
}
