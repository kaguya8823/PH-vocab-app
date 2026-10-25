package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"PH-vocab-app/backend/handler"
	"PH-vocab-app/backend/repository"
	"PH-vocab-app/backend/router"
	"PH-vocab-app/backend/usecase"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN must be set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open database connection: %v", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	if err := repository.EnsureSchema(db); err != nil {
		log.Fatalf("initialize database schema: %v", err)
	}

	vocabularyRepository := repository.NewVocabularyRepository(db)
	vocabularyUseCase := usecase.NewVocabularyUseCase(vocabularyRepository)
	vocabularyHandler := handler.NewVocabularyHandler(vocabularyUseCase)

	// 音声用のRepositoryを作成
	audioRepository := repository.NewAudioRepository(db)

	// 音声用のUseCaseを作成
	audioUseCase := usecase.NewAudioUseCase(audioRepository)

	// 音声用のHandlerを作成
	audioHandler := handler.NewAudioHandler(audioUseCase)

	server := &http.Server{
		Addr: ":3002",
		Handler: router.NewVocabularyRouter(
			vocabularyHandler,
			audioHandler,
		),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve HTTP: %v", err)
	}
}
