package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"PH-vocab-app/backend/repository"
	"PH-vocab-app/backend/usecase"
)

type AudioHandler struct {
	audioUseCase usecase.AudioUseCase
}

func NeaAudioHandler(audioUseCase usecase.AudioUseCase) *AudioHandler {
	return &AudioHandler{
		audioUseCase: audioUseCase,
	}
}

// URLから単語IDを取得
// 例：　/api/vocabulary/12/audio -> "12"
func getVocabularyIDFromAudioPath(r *http.Request) string {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 4 ||
		parts[0] != "api" ||
		parts[1] != "vocabulary" ||
		parts[3] != "audio" {
		return ""
	}

	return parts[2]
}

// POST /api/vocabulary/{id}/audio
// 音声の保存情報を登録・更新
func (h *AudioHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vocabularyID := getVocabularyIDFromAudioPath(r)
	if vocabularyID == "" {
		http.Error(w, "invalid vocabulary ID path", http.StatusBadRequest)
		return
	}

	var request struct {
		AudioType string `json:"audio_type"`
		FilePath  string `json:"file_path"`
		MimeType  string `json:"mime_type"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // Jsonは最大1MB
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	audio := repository.Audio{
		VocabularyID: vocabularyID,
		AudioType:    request.AudioType,
		FilePath:     request.FilePath,
		MimeType:     request.MimeType,
	}

	if err := h.audioUseCase.Save(audio); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "audio metadata saved",
	})
}

// GET /api/vocabulary/{id}/audio?type=word
// 音声の保存情報を取得
func (h *AudioHandler) Find(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vocabularyID := getVocabularyIDFromAudioPath(r)
	if vocabularyID == "" {
		http.Error(w, "invalid vocabulary ID path", http.StatusBadRequest)
		return
	}

	audioType := r.URL.Query().Get("type")

	audio, err := h.audioUseCase.FindByVocabularyIDAndType(
		vocabularyID,
		audioType,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if audio == nil {
		http.Error(w, "audio not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(audio); err != nil {
		// 応答開始後のエラーでは、HTTPステータスを変更できない
		return
	}
}
