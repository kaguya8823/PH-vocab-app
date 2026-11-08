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
    id: string,
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

export const deleteVocabulary = async (id: string): Promise<void> => {
    await axios.delete(`${API_URL}/${id}`);
};

// 音声API処理
export type AudioType = 'word' | 'example';

// 音声を Go API にアップロード
export const uploadVocabularyAudio = async (
    vocabularyId: string,
    audioType:    AudioType,
    audioBlob:    Blob,
): Promise<void> => {
    const formData = new FormData();

    // Go API が受け取るフィールド名に合わせる
    formData.append('audio_type', audioType);

    const extension = audioBlob.type.includes('ogg')
        ? 'ogg'
        : audioBlob.type.includes('mp4')
            ? 'mp4'
            : audioBlob.type.includes('wav')
            ? 'wav'
            : 'webm';

    formData.append(
        'audio',
        audioBlob,
        `${audioType}.${extension}`,
    );

    await axios.post(
        `/api/vocabulary/${vocabularyId}/audio`,
        formData,
    );
};

// 保存された音声を再生するための URL 取得
export const getVocabularyAudioUrl = (
    vocabularyId: string,
    audioType:    AudioType,
): string => {
    return `/api/vocabulary/${vocabularyId}/audio?type=${audioType}`;
}

// お気に入り音声のメタデータ定義
export interface AudioMetadata {
    vocabularyId: string;
    audio_type:   AudioType;
    mime_type:    string;
    is_favorite:  boolean;
}

// お気に入り音声のメタデータを取得
export const getVocabularyAudioMetadata = async (
    vocabularyId: string,
    audioType:    AudioType,
): Promise<AudioMetadata> => {
    const response = await axios.get<AudioMetadata>(
        `/api/vocabulary/${vocabularyId}/audio/metadata`,
        {
            params: {
                type: audioType,
            },
        },
    );

    return response.data;
};

// お気に入り音声のメタデータを更新
export const setVocabularyAudioFavorite = async (
    vocabularyId: string,
    audioType:    AudioType,
    isFavorite:   boolean,
): Promise<void> => {
    await axios.put(
        `/api/vocabulary/${vocabularyId}/audio/favorite`,
        {
            audio_type:  audioType,
            is_favorite: isFavorite,
        },
    );
};