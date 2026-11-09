import {useEffect, useState} from 'react';
import { isAxiosError } from 'axios';

import {
    getVocabularyAudioMetadata,
    setVocabularyAudioFavorite,
    deleteVocabularyAudio
} from '../services/vocabulary.service'; 

import type {
    AudioMetadata,
    AudioType
} from '../services/vocabulary.service';

type AudioItemProps = {
    vocabularyId: string;
    audioType:    AudioType;
    label:        string;
};

type Props = {
    vocabularyId: string;
};

const AudioItem = ({
    vocabularyId,
    audioType,
    label,
}: AudioItemProps) => {
    const [metadata, setMetadata] = useState<AudioMetadata | null>(null);
    const [loading, setLoading] = useState(true);
    const [processing, setProcessing] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const loadMetadata = async () => {
        setLoading(true);
        setError(null);
        try {
            const result = await getVocabularyAudioMetadata(
                vocabularyId,
                audioType
            );
            setMetadata(result);
        } catch (err) {
            if (isAxiosError(err) && err.response?.status === 404) {
                setMetadata(null);
            } else {
                setError("音声情報を取得できませんでした。");
            }
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        let cancelled = false;

        const load = async () => {
            try {
                const result = await getVocabularyAudioMetadata(
                    vocabularyId,
                    audioType,
                );

                if (!cancelled) {
                    setMetadata(result);
                }
            } catch (err) {
                if (!cancelled) {
                    if (isAxiosError(err) && err.response?.status === 404) {
                        setMetadata(null);
                    } else {
                        setError("音声情報を取得できませんでした。");
                    }
                }
            } finally {
                if (!cancelled) {
                    setLoading(false);
                }
            }
        };

        void load();

        return () => {
            cancelled = true;
        };
    }, [vocabularyId, audioType]);

    const handleFavorite = async () => {
        if (!metadata) return;

        setProcessing(true);
        setError(null);

        try {
            const nextFavorite = !metadata.is_favorite;

            await setVocabularyAudioFavorite(
                vocabularyId,
                audioType,
                nextFavorite,
            );

            setMetadata({
                ...metadata,
                is_favorite: nextFavorite,
            });
        } catch {
            setError("お気に入り状態の更新に失敗しました。");
        } finally {
            setProcessing(false);
        }
    };

    const handleDelet = async () => {
        if (!metadata) return;

        if (metadata.is_favorite) {
            setError("お気に入り音声は削除できません。先に解除してください。");
            return;
        }

        const confirmed = window.confirm(
            `${label} の音声を削除しますか？取り消しはできません。`
        );

        if (!confirmed) return;

        setProcessing(true);
        setError(null);

        try {
            await deleteVocabularyAudio(
                vocabularyId,
                audioType,
            );
            setMetadata(null);
        } catch (err) {
            if (isAxiosError(err) && err.response?.status === 409) {
                setError("お気に入り音声は削除できません。");
            } else {
                setError("音声の削除に失敗しました。");
            }
        } finally {
            setProcessing(false);
        }
    };

    if (loading) {
        return <p>{label}の情報を読み込み中...</p>;
    }

    return (
        <section
            style={{
                border: "1px solid #aaa",
                borderRadius: "6px",
                padding: "10px",
                marginBottom: "10px",
            }}
        >
            <h4 style={{ marginTop: 0 }}>{label}</h4>

            {metadata ? (
                <>
                    <p>
                        状態：
                        {metadata.is_favorite
                            ? "お気に入り"
                            : "通常"}
                    </p>

                    <button
                        type="button"
                        onClick={handleFavorite}
                        disabled={processing}
                    >
                        {metadata.is_favorite
                            ? "お気に入り解除"
                            : "お気に入り登録"}
                    </button>

                    <button
                        type="button"
                        onClick={handleDelet}
                        disabled={processing || metadata.is_favorite}
                        style={{ marginLeft: "10px" }}
                    >
                        削除
                    </button>
                </>
            ) : (
                <p>音声は登録されていません。</p>
            )}

            {processing && <p>処理中...</p>}

            {error && (
                <p role='alert' style={{ color: "red"}}>
                    {error}
                </p>
            )}

            {!loading && (
                <button
                    type='button'
                    onClick={() => void loadMetadata()}
                    disabled={processing}
                    style={{ marginTop: "8px" }}
                >
                    情報を読み込み
                </button>
            )}
        </section>
    );
};

const VocabularyAudioManagement = ({ vocabularyId }: Props) => {
    return (
        <div>
            <h3>音声管理</h3>
            <AudioItem
                vocabularyId={vocabularyId}
                audioType="word"
                label="単語音声"
            />
            <AudioItem
                vocabularyId={vocabularyId}
                audioType="example"
                label="例文音声"
            />
        </div>
    );
};

export default VocabularyAudioManagement;
