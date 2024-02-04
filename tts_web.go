package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"tts-web/db"
	"tts-web/model"
	"tts-web/service"
)

func main() {
	database := db.InitializeDB()
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("Error closing the database connection: %v", err)
		}
	}(database)

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/tts", textToSpeechHandler)
	r.HandleFunc("/api/tts", handleTTSRequest(database))
	r.HandleFunc("/api/tts/records", handleTTSRecordsRequest(database))
	r.HandleFunc("/api/tts/records/{recordId}/audio", handleTTSRecordAudioByID(database))
	r.HandleFunc("/api/tts/records/{recordId}", handleDeleteTTSRecordByID(database))

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	port := os.Getenv("TTS_PORT")
	if port == "" {
		port = "8081"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func textToSpeechHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tts" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "This is the TTS page.")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleTTSRequest(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ttsService := service.OpenAITTSService{}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Text  string `json:"text"`
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ttsReq := model.TTSRequest{
			Model: "tts-1-hd",
			Input: request.Text,
			Voice: "nova",
		}

		audio, err := ttsService.GenerateSpeech(ttsReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.SaveTTSAudio(database, model.TTSRecord{
			Title:        request.Title,
			InputText:    request.Text,
			AudioContent: audio,
		})
		if err != nil {
			log.Fatalf("Error saving TTS audio: %v", err)
		}
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write(audio)
	}
}

func handleTTSRecordAudioByID(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		vars := mux.Vars(r)
		recordID := vars["recordId"]
		audio, err := db.FetchTTSAudioByID(database, recordID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write(audio)
	}
}

func handleDeleteTTSRecordByID(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		vars := mux.Vars(r)
		recordID := vars["recordId"]
		err := db.DeleteTTSRecordByID(database, recordID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleTTSRecordsRequest(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		records, err := db.FetchAllTTSRecords(database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse, err := json.Marshal(records)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
