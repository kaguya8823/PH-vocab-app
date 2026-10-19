package domain

type Vocabulary struct {
	ID              string `json:"id"`
	Word            string `json:"word"`
	Meaning         string `json:"meaning"`
	ExampleSentence string `json:"exampleSentence"`
	PartOfSpeech    string `json:"partOfSpeech"`
	Comment         string `json:"comment"`
	CreatedAt       string `json:"createdAt"`
}

type CreateVocabularyInput struct {
	Word            string `json:"word"`
	Meaning         string `json:"meaning"`
	ExampleSentence string `json:"exampleSentence"`
	PartOfSpeech    string `json:"partOfSpeech"`
	Comment         string `json:"comment"`
}

// ここにはサンプルボイスを入れたい
