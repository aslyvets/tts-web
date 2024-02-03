package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"tts-web/model"
	"tts-web/service"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/tts", textToSpeechHandler)
	http.HandleFunc("/api/tts", handleTTSRequest)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := os.Getenv("TTS_PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func textToSpeechHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tts" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
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

func handleTTSRequest(w http.ResponseWriter, r *http.Request) {
	ttsService := service.OpenAITTSService{}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Text string `json:"text"`
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

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(audio)
}
