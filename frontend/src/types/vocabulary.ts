export type Vocabulary = {
    id:               number;
    word:             string;
    meaning:          string;
    exampleSentence?: string;
    paerOfSpeech:     string;

    createdAt:        string;
    updatedAt:        string;
}

// ここには録音した音声を付けたい