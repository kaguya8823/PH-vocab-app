package usecase

import (
	"PH-vocab-app/backend/domain"
	"PH-vocab-app/backend/repository"
)

type VocabularyUseCase struct {
	repository VocabularyRepository
}

func (u *VocabularyUseCase) GetVocabularies()
	([]domain.Vocabulary, error) {

		return u.repository.FindAll()
	}