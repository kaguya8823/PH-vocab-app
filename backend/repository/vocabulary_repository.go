package repository

import (
	"PH-vocab-app/backend/domain"
	"database/sql"
	"strconv"
)

type VocabularyRepository interface {
	FindAll() ([]domain.Vocabulary, error)
	Create(domain.Vocabulary) (domain.Vocabulary, error)
	Update(string, domain.Vocabulary) (bool, error)
	Delete(string) (bool, error)
}

type vocabularyRepository struct {
	db *sql.DB
}

func NewVocabularyRepository(db *sql.DB) VocabularyRepository {
	return &vocabularyRepository{db: db}
}

func EnsureSchema(db *sql.DB) error {
	// 既存の単語テーブル
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
	if err != nil {
		return err
	}

	// 新しい音声テーブル
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vocabulary_audios (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			vocabulary_id BIGINT UNSIGNED NOT NULL,
			audio_type VARCHAR(20) NOT NULL,
			file_path VARCHAR(500) NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			UNIQUE KEY uq_vocabulary_audio_type
				(vocabulary_id, audio_type),

			CONSTRAINT fk_vocabulary_audios_vocabulary
				FOREIGN KEY (vocabulary_id)
				REFERENCES vocabularies(id)
				ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	if err != nil {
		return err
	}

	return nil
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

func (r *vocabularyRepository) Update(id string, vocabulary domain.Vocabulary) (bool, error) {
	result, err := r.db.Exec(`
		UPDATE vocabularies
		SET word = ?, meaning = ?, example_sentence = ?, part_of_speech = ?, comment = ?
		WHERE id = ?
	`,
		vocabulary.Word,
		vocabulary.Meaning,
		vocabulary.ExampleSentence,
		vocabulary.PartOfSpeech,
		vocabulary.Comment,
		id,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected > 0 {
		return true, nil
	}

	var exists bool
	if err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM vocabularies WHERE id = ?)", id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *vocabularyRepository) Delete(id string) (bool, error) {
	result, err := r.db.Exec("DELETE FROM vocabularies WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}
