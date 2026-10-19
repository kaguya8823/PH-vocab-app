package router

import (
	"net/http"

	"PH-vocab-app/backend/handler"
	"github.com/gorilla/mux"
)

func NewVocabularyRouter(vocabularyHandler *handler.VocabularyHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(
		"/api/vocabularies",
		vocabularyHandler.GetVocabularies,
	).Methods(http.MethodGet)
	return router
}
