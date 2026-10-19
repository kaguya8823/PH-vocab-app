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