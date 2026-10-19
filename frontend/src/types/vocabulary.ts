export type Vocabulary = {
    id:               string;
    word:             string;
    meaning:          string;
    exampleSentence?: string;
    partOfSpeech:     string;
    comment:          string;
    createdAt:        string;
};

export type CreateVocabularyInput = {
    word: string;
    meaning: string;
    exampleSentence?: string;
    partOfSpeech: string;
    comment: string;
};