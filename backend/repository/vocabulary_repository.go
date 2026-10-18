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

func (r *vocabularyRepository) FindAll() ([]domain.Vocabulary, error) {

	rows, err := r.db.Query(`
	SELECT
		id,
		word,
		meaning,
		example_sentence,
		part_of_speech,
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

		err := rows.Scan(
			&vocabulary.ID,
			&vocabulary.Word,
			&vocabulary.Meaning,
			&vocabulary.ExampleSentence,
			&vocabulary.PartOfSpeech,
		)

		if err != nil {
			return nil, err
		}

		vocabularies =
			append(vocabularies, vocabulary)
	}

	return vocabularies, nil
}
