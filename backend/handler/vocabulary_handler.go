package handler

import (
	"encoding/json"
	"net/http"

	"PH-vocab-app/backend/domain"
)

type vocabularyUseCase interface {
	GetVocabularies() ([]domain.Vocabulary, error)
}

type VocabularyHandler struct {
	useCase vocabularyUseCase
}

func NewVocabularyHandler(useCase vocabularyUseCase) *VocabularyHandler {
	return &VocabularyHandler{useCase: useCase}
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
