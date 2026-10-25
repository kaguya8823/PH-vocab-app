package usecase

import (
	"fmt"
	"strconv"
	"strings"

	"PH-vocab-app/backend/repository"
)

type AudioUseCase interface {
	Save(audio repository.Audio) error
	FindByVocabularyIDAndType(
		vocabularyID string,
		audioType string,
	) (*repository.Audio, error)
}

type audioUseCase struct {
	audioRepository repository.AudioRepository
}

func NewAudioUseCase(
	audioRepository repository.AudioRepository,
) AudioUseCase {
	return &audioUseCase{
		audioRepository: audioRepository,
	}
}

// 音声情報を検証して保存
func (u *audioUseCase) Save(audio repository.Audio) error {
	// 単語IDが正しいか確認
	vocabularyID := strings.TrimSpace(audio.VocabularyID)
	id, err := strconv.ParseUint(vocabularyID, 10, 64)
	if err != nil || id == 0 {
		return fmt.Errorf("invalid vocabulary ID")
	}
	audio.VocabularyID = vocabularyID

	// 音声の種類を確認
	switch audio.AudioType {
	case "word", "example":
		// 許可された種類なら続行
	default:
		return fmt.Errorf("invalid audio type: must be word or example")
	}

	// 保存先と音声形式が指定されているか確認
	if strings.TrimSpace(audio.FilePath) == "" {
		return fmt.Errorf("file path is required")
	}
	if strings.TrimSpace(audio.MimeType) == "" {
		return fmt.Errorf("mime type is required")
	}

	// 検証に成功したら Repository へ登録依頼する
	return u.audioRepository.Save(audio)
}

// 単語IDと音声の種類から音声情報を取得
func (u *audioUseCase) FindByVocabularyIDAndType(
	vocabularyID string,
	audioType string,
) (*repository.Audio, error) {
	vocabularyID = strings.TrimSpace(vocabularyID)

	id, err := strconv.ParseUint(vocabularyID, 10, 64)
	if err != nil || id == 0 {
		return nil, fmt.Errorf("invalid vocabulary ID")
	}

	switch audioType {
	case "word", "example":
		// 許可された種類なので続行
	default:
		return nil, fmt.Errorf("invalid audio type: must be word or example")
	}

	return u.audioRepository.FindByVocabularyIDAndType(
		vocabularyID,
		audioType,
	)
}
