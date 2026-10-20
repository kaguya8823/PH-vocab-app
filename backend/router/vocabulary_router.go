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
	router.HandleFunc(
		"/api/vocabularies",
		vocabularyHandler.CreateVocabulary,
	).Methods(http.MethodPost)
	router.HandleFunc(
		"/api/vocabularies/{id}",
		vocabularyHandler.UpdateVocabulary,
	).Methods(http.MethodPut)
	router.HandleFunc(
		"/api/vocabularies/{id}",
		vocabularyHandler.DeleteVocabulary,
	).Methods(http.MethodDelete)
	return router
}
