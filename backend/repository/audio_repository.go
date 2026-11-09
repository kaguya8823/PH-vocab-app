package repository

import (
	"database/sql"
	"errors"
)

type Audio struct {
	ID           int64  `json:"id"`
	VocabularyID string `json:"vocabulary_id"`
	AudioType    string `json:"audio_type"`
	FilePath     string `json:"file_path"`
	MimeType     string `json:"mime_type"`
	CreatedAt    string `json:"created_at"`
	IsFavorite   bool   `json:"is_favorite"`
}

// お気に入りは削除できない定義
var ErrAudioIsFavorite = errors.New(
	"cannot delete favorite audio",
)

type AudioRepository interface {
	Save(audio Audio) error
	FindByVocabularyIDAndType(
		vocabularyID string,
		audioType string,
	) (*Audio, error)
	SetFavorite(
		vocabularyID string,
		audioType string,
		isFavorite bool,
	) error
	// 削除処理
	Delete(
		vocabularyID string,
		audioType string,
	) (*Audio, error)
}

type audioRepository struct {
	db *sql.DB
}

func NewAudioRepository(db *sql.DB) AudioRepository {
	return &audioRepository{db: db}
}

// 音声保存の登録
// 同じ単語、種類の音声がすでにあれば更新
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
			file_path = IF(is_favorite = 0, VALUES(file_path), file_path),
			mime_type = IF(is_favorite = 0, VALUES(mime_type), mime_type)
	`,
		audio.VocabularyID,
		audio.AudioType,
		audio.FilePath,
		audio.MimeType,
	)
	if err != nil {
		return err
	}
	return nil
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
			CAST(created_at AS CHAR),
			is_favorite
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
		&audio.IsFavorite,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &audio, nil
}

// お気に入り状態を更新する
func (r *audioRepository) SetFavorite(
	vocabularyID string,
	audioType string,
	isFavorite bool,
) error {
	result, err := r.db.Exec(`
		UPDATE vocabulary_audios
		SET is_favorite = ?
		WHERE vocabulary_id = ? AND audio_type = ?
	`,
		isFavorite,
		vocabularyID,
		audioType,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// 値が変わらない場合、RowsAffectedが０になることがある。
	// その場合はエラーにせず、正常終了とする。
	if rowsAffected == 0 {
		var exists bool
		err = r.db.QueryRow(`
			SELECT EXISTS(
				SELECT 1
				FROM vocabulary_audios
				WHERE vocabulary_id = ? AND audio_type = ?
			)
		`, vocabularyID, audioType).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return sql.ErrNoRows
		}
	}
	return nil
}

// お気に入り保護を確認、音声情報を削除する
func (r *audioRepository) Delete(
	vocabularyID string,
	audioType string,
) (*Audio, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var audio Audio

	// 削除対象の音声路ロックして取得する
	err = tx.QueryRow(`
	SELECT
		id,
		vocabulary_id,
		audio_type,
		file_path,
		mime_type,
		CAST(created_at AS CHAR),
		is_favorite
	FROM vocabulary_audios
	WHERE vocabulary_id = ?
		AND audio_type = ?
	FOR UPDATE
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
		&audio.IsFavorite,
	)

	if err != nil {
		return nil, err
	}

	// お気に入りは削除できない
	if audio.IsFavorite {
		return nil, ErrAudioIsFavorite
	}

	_, err = tx.Exec(`
		DELETE FROM vocabulary_audios
		WHERE vocabulary_id = ?
		 AND audio_type = ?
	`,
		vocabularyID,
		audioType,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &audio, nil
}
