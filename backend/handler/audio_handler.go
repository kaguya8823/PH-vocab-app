package handler

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"PH-vocab-app/backend/repository"
	"PH-vocab-app/backend/usecase"
)

const maxUploadSize int64 = 10 << 20 // 10Mib

type AudioHandler struct {
	audioUseCase usecase.AudioUseCase
	audioStorage *usecase.AudioStorage
}

func NewAudioHandler(
	audioUseCase usecase.AudioUseCase,
	audioStorage *usecase.AudioStorage,
) *AudioHandler {
	return &AudioHandler{
		audioUseCase: audioUseCase,
		audioStorage: audioStorage,
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

// 音声ファイルの形式
func audioFormat(contentType string) (string, string, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", "", fmt.Errorf("invalid audio content type")
	}

	switch strings.ToLower(mediaType) {
	case "audio/webm":
		return ".webm", "audio/webm", nil
	case "audio/ogg":
		return ".ogg", "audio/ogg", nil
	case "audio/mp4":
		return ".mp4", "audio/mp4", nil
	case "audio/wav", "audio/x-wav", "audio/wave":
		return ".wav", "audio/wav", nil
	default:
		return "", "", fmt.Errorf("unsupported audio format")
	}
}

// Save reveives an audio file and stores its metadata.
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

	// Allow a small amount of multipart overhead.
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+(1<<20))
	defer r.Body.Close()

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		http.Error(w, "invalid or oversized multipart request", http.StatusBadRequest)
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	audioType := r.FormValue("audio_type")
	if audioType != "word" && audioType != "example" {
		http.Error(w, "audio_type must be word or example", http.StatusBadRequest)
		return
	}

	audioFile, fileHeader, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "audio file is required", http.StatusBadRequest)
		return
	}
	defer audioFile.Close()

	if fileHeader.Size <= 0 || fileHeader.Size > maxUploadSize {
		http.Error(w, "audio file must be between 1 byte and 10 MiB", http.StatusBadRequest)
		return
	}

	extension, contentType, err := audioFormat(
		fileHeader.Header.Get("Content-Type"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Look up the previous recording so it can be removed after replavement.
	oldAudio, err := h.audioUseCase.FindByVocabularyIDAndType(
		vocabularyID,
		audioType,
	)
	if err != nil {
		http.Error(w, "failed to look up existing audio", http.StatusInternalServerError)
		return
	}

	filePath, err := h.audioStorage.Save(audioFile, extension)
	if err != nil {
		http.Error(w, "failed to save audio file", http.StatusInternalServerError)
		return
	}

	audio := repository.Audio{
		VocabularyID: vocabularyID,
		AudioType:    audioType,
		FilePath:     filePath,
		MimeType:     contentType,
	}

	if err := h.audioUseCase.Save(audio); err != nil {
		_ = h.audioStorage.Remove(filePath)
		http.Error(w, "failed to save audio metadata", http.StatusInternalServerError)
		return
	}

	// Remoce the previous file only after the new metadata has been saved.
	if oldAudio != nil && oldAudio.FilePath != filePath {
		if err := h.audioStorage.Remove(oldAudio.FilePath); err != nil {
			// The new recording is saved successfully; report the cleanup issue
			// in server logs rather than treating the upload as failed.
			fmt.Fprintf(os.Stderr, "remove old audio file: %v\n", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "audio uploaded successfully",
	})
}

// Find returns audio metadata as JSON.
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
	if audioType != "word" && audioType != "example" {
		http.Error(w, "type must be word or example", http.StatusBadRequest)
		return
	}

	// MySQL から音声情報を取得する
	audio, err := h.audioUseCase.FindByVocabularyIDAndType(
		vocabularyID,
		audioType,
	)
	if err != nil {
		http.Error(w, "filed to retrieve audio information", http.StatusInternalServerError)
		return
	}

	if audio == nil {
		http.Error(w, "audio not found", http.StatusNotFound)
		return
	}

	// 保存済みの音声ファイルを開く
	file, info, err := h.audioStorage.Open(audio.FilePath)
	if err != nil {
		http.Error(w, "audio file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", audio.MimeType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")

	// 音声データを返す。ブラウザの再位置指定にも対応
	http.ServeContent(
		w,
		r,
		filepath.Base(audio.FilePath),
		info.ModTime(),
		file,
	)
}

// Keep filepath imported only if you use it for additional path checks later.
var _ = filepath.Base
