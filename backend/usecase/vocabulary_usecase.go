package usecase

import (
	"PH-vocab-app/backend/domain"
	"PH-vocab-app/backend/repository"
)

type VocabularyUseCase struct {
	repository repository.VocabularyRepository
}

func NewVocabularyUseCase(repository repository.VocabularyRepository) *VocabularyUseCase {
	return &VocabularyUseCase{repository: repository}
}

func (u *VocabularyUseCase) GetVocabularies() ([]domain.Vocabulary, error) {
	return u.repository.FindAll()
}

func (u *VocabularyUseCase) CreateVocabulary(input domain.CreateVocabularyInput) (domain.Vocabulary, error) {
	return u.repository.Create(domain.Vocabulary{
		Word:            input.Word,
		Meaning:         input.Meaning,
		ExampleSentence: input.ExampleSentence,
		PartOfSpeech:    input.PartOfSpeech,
		Comment:         input.Comment,
	})
}

func (u *VocabularyUseCase) UpdateVocabulary(id string, input domain.UpdateVocabularyInput) (bool, error) {
	return u.repository.Update(id, domain.Vocabulary{
		ID:              id,
		Word:            input.Word,
		Meaning:         input.Meaning,
		ExampleSentence: input.ExampleSentence,
		PartOfSpeech:    input.PartOfSpeech,
		Comment:         input.Comment,
	})
}
