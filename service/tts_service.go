package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"tts-web/model"
)

type TTSService interface {
	GenerateSpeech(ttsReq model.TTSRequest) ([]byte, error)
}

type OpenAITTSService struct {
}

func (s *OpenAITTSService) GenerateSpeech(ttsReq model.TTSRequest) ([]byte, error) {
	requestBody, err := json.Marshal(ttsReq)
	apiKey := os.Getenv("TTS_API_KEY")
	if err != nil {
		return nil, err
	}

	openAIURL := "https://api.openai.com/v1/audio/speech"
	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
