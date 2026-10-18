import axios from 'axios';

const API_URL =
    'http://localhost:3003/api/vocabularies';

export const fetchVocabularies = async () => {

    const response =
        await axios.get(API_URL);
    
    return response.data;
};

export const createVocabulary = async (
    data: {
        word:             string;
        meaning:          string;
        exampleSentence?: string;
        partOfSpeech:     string;
    }
) => {

    const response =
        await axios.post(API_URL, data);

    return response.data;
};