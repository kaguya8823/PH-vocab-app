package repository

import (
	"PH-vocab-app/backend/domain"
	"database/sql"
)

type VocabularyRepository interface {
	FindAll() ([]domain.Vocabulary, error)
}

type vocabularyRepository struct {
	db *sql.DB
}

func NewVocabularyRepository(db *sql.DB) VocabularyRepository {
	return &vocabularyRepository{db: db}
}

func (r *vocabularyRepository) FindAll() ([]domain.Vocabulary, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			word,
			meaning,
			example_sentence,
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

	var vocabularies []domain.Vocabulary
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
