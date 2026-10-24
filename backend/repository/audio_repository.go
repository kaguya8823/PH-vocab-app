package repository

import (
	"database/sql"
)

type Audio struct {
	ID           int64  `json:"id"`
	VocabularyID string `json:"vocabulary_id"`
	AudioType    string `json:"audio_type"`
	FilePath     string `json:"file_path"`
	MimeType     string `json:"mime_type"`
	CreatedAt    string `json:"created_at"`
}

type AudioRepository interface {
	Save(audio Audio) error
	FindByVocabularyIDAndType(
		vocabularyID string,
		audioType string,
	) (*Audio, error)
}

type audioRepository struct {
	db *sql.DB
}

func NeaAudioRepository(db *sql.DB) AudioRepository {
	return &audioRepository{db: db}
}

// 音声保存の登録
// 同じ単語、酒類の音声がすでにあれば更新
func (r *audioRepository) Save(audio Audio) error {
	_, err := r.db.Exec(`
		INSERT INTO vocabulary_audios (
			vocabulary_id,
			audio_type,
			file_path,
			mime_type
		)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			file_path = VALUES(file_path),
			mime_type = VALUES(mime_type)
	`,
		audio.VocabularyID,
		audio.AudioType,
		audio.FilePath,
		audio.MimeType,
	)

	return err
}

// 単語のIDと音声の種類から保存情報を取得
func (r *audioRepository) FindByVocabularyIDAndType(
	vocabularyID string,
	audioType string,
) (*Audio, error) {
	var audio Audio

	err := r.db.QueryRow(`
		SELECT
			id,
			vocabulary_id,
			audio_type,
			file_path,
			mime_type,
			CAST(created_at AS CHAR)
		FROM vocabulary_audios
		WHERE vocabulary_id = ?
			AND audio_type = ?
	`,
		vocabularyID,
		audioType,
	).Scan(
		&audio.ID,
		&audio.VocabularyID,
		&audio.AudioType,
		&audio.FilePath,
		&audio.MimeType,
		&audio.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &audio, nil
}
