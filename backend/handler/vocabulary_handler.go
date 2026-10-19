package handler

import (
	"encoding/json"
	"net/http"

	"PH-vocab-app/backend/domain"
)

type vocabularyUseCase interface {
	GetVocabularies() ([]domain.Vocabulary, error)
	CreateVocabulary(domain.CreateVocabularyInput) (domain.Vocabulary, error)
}

type VocabularyHandler struct {
	useCase vocabularyUseCase
}

func NewVocabularyHandler(useCase vocabularyUseCase) *VocabularyHandler {
	return &VocabularyHandler{useCase: useCase}
}

func (h *VocabularyHandler) CreateVocabulary(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input domain.CreateVocabularyInput
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	if input.Word == "" || input.Meaning == "" || input.PartOfSpeech == "" {
		http.Error(w, "word, meaning, and partOfSpeech are required", http.StatusBadRequest)
		return
	}

	vocabulary, err := h.useCase.CreateVocabulary(input)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(vocabulary); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *VocabularyHandler) GetVocabularies(
	w http.ResponseWriter,
	r *http.Request,
) {
	vocabularies, err :=
		h.useCase.GetVocabularies()

	if err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(vocabularies); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
