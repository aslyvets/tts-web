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
