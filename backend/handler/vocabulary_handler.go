package handler

import "net/http"

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

	w.Handler().Set(
		"Content-Type",
		"application/json",
	)

	json.NetEncoder(w).Encode(vocabularies)
}
