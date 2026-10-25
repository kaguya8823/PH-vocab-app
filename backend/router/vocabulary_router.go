package router

import (
	"net/http"

	"PH-vocab-app/backend/handler"

	"github.com/gorilla/mux"
)

func NewVocabularyRouter(
	vocabularyHandler *handler.VocabularyHandler,
	audioHandler *handler.AudioHandler,
) *mux.Router {
	router := mux.NewRouter()

	// 単語一覧の取得
	router.HandleFunc(
		"/api/vocabularies",
		vocabularyHandler.GetVocabularies,
	).Methods(http.MethodGet)

	// 単語の登録
	router.HandleFunc(
		"/api/vocabularies",
		vocabularyHandler.CreateVocabulary,
	).Methods(http.MethodPost)

	// 単語の更新
	router.HandleFunc(
		"/api/vocabularies/{id}",
		vocabularyHandler.UpdateVocabulary,
	).Methods(http.MethodPut)

	// 単語の削除
	router.HandleFunc(
		"/api/vocabularies/{id}",
		vocabularyHandler.DeleteVocabulary,
	).Methods(http.MethodDelete)

	// 音声情報の保存
	router.HandleFunc(
		"/api/vocabulary/{id}/audio",
		audioHandler.Save,
	).Methods(http.MethodPost)

	// 音声情報の取得
	router.HandleFunc(
		"/api/vocabulary/{id}/audio",
		audioHandler.Find,
	).Methods(http.MethodGet)

	return router
}
