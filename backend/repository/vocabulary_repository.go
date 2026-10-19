package repository

import (
	"PH-vocab-app/backend/domain"
	"database/sql"
	"strconv"
)

type VocabularyRepository interface {
	FindAll() ([]domain.Vocabulary, error)
	Create(domain.Vocabulary) (domain.Vocabulary, error)
}

type vocabularyRepository struct {
	db *sql.DB
}

func NewVocabularyRepository(db *sql.DB) VocabularyRepository {
	return &vocabularyRepository{db: db}
}

func EnsureSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS vocabularies (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			word VARCHAR(255) NOT NULL,
			meaning TEXT NOT NULL,
			example_sentence TEXT NOT NULL,
			part_of_speech VARCHAR(50) NOT NULL,
			comment TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	return err
}

func (r *vocabularyRepository) FindAll() ([]domain.Vocabulary, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			word,
			meaning,
			COALESCE(example_sentence, ''),
			part_of_speech,
			COALESCE(comment, ''),
			CAST(created_at AS CHAR)
		FROM vocabularies
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vocabularies := make([]domain.Vocabulary, 0)
	for rows.Next() {
		var vocabulary domain.Vocabulary
		if err := rows.Scan(
			&vocabulary.ID,
			&vocabulary.Word,
			&vocabulary.Meaning,
			&vocabulary.ExampleSentence,
			&vocabulary.PartOfSpeech,
			&vocabulary.Comment,
			&vocabulary.CreatedAt,
		); err != nil {
			return nil, err
		}

		vocabularies = append(vocabularies, vocabulary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return vocabularies, nil
}

func (r *vocabularyRepository) Create(vocabulary domain.Vocabulary) (domain.Vocabulary, error) {
	result, err := r.db.Exec(`
		INSERT INTO vocabularies (
			word,
			meaning,
			example_sentence,
			part_of_speech,
			comment
		) VALUES (?, ?, ?, ?, ?)
	`,
		vocabulary.Word,
		vocabulary.Meaning,
		vocabulary.ExampleSentence,
		vocabulary.PartOfSpeech,
		vocabulary.Comment,
	)
	if err != nil {
		return domain.Vocabulary{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Vocabulary{}, err
	}
	vocabulary.ID = strconv.FormatInt(id, 10)
	return vocabulary, nil
}
