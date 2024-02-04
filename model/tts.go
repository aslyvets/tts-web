package model

type TTSRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
}

type TTSResponse struct {
	AudioContent []byte
}

type TTSRecord struct {
	Title        string
	InputText    string
	AudioContent []byte
}

type TTSListRecord struct {
	Id    string
	Title string
	Text  string
}
