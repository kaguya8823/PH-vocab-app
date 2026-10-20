import axios from 'axios';
import type { CreateVocabularyInput, Vocabulary } from '../types/vocabulary';

const API_URL = '/api/vocabularies';

export const fetchVocabularies = async (): Promise<Vocabulary[]> => {
    const response = await axios.get<Vocabulary[]>(API_URL);
    return response.data;
};

export const createVocabulary = async (
    data: CreateVocabularyInput,
): Promise<Vocabulary> => {
    const response = await axios.post<Vocabulary>(API_URL, data);
    return response.data;
};

export const updateVocabulary = async (
    id: number,
    data: {
        word:             string;
        meaning:          string;
        exampleSentence?: string;
        partOfSpeech:     string;
        comment:          string;
    }
) => {
    const response = await axios.put(
        `${API_URL}/${id}`,
        data
    );

    return response.data;
}

export const deleteVocabulary = async (id: number) => {
    const response = await axios.delete(`${API_URL}/${id}`);
    return response.data;
}