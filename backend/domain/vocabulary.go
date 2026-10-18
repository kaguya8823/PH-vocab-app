package domain

type Vocabulary struct {
	ID              string `json:"id"`
	Word            string `json:"word"`
	Meaning         string `json:"meaning"`
	ExampleSentence string `json:"exampleSentence"`
	PartOfSpeech    string `json:"partOfSpeech"`
}

// ここにはサンプルボイスを入れたい
