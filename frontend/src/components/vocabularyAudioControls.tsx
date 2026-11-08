import { useEffect, useRef, useState } from "react";
import {
    getVocabularyAudioUrl,
    getVocabularyAudioMetadata,
    setVocabularyAudioFavorite,
    uploadVocabularyAudio,
} from '../services/vocabulary.service';
import type { AudioType } from '../services/vocabulary.service';
import { isAxiosError } from 'axios';

type Props = {
    vocabularyId: string;
    audioType:    AudioType;
    label:        string;
};

const VocabularyAudioControls = ({
    vocabularyId,
    audioType,
    label,
}: Props) => {
    const recorderRef = useRef<MediaRecorder | null>(null);
    const streamRef = useRef<MediaStream | null>(null);
    const chunksRef = useRef<Blob[]>([]);

    const [recording, setRecording] = useState(false);
    const [saving, setSaving] = useState(false);
    const [audioVersion, setAudioVersion] = useState(0);
    const [error, setError] = useState('');

    const [isFavorite, setIsFavorite] = useState(false);
    const [favoriteLoading, setFavoriteLoading] = useState(true);

    useEffect(() => {
        return () => {
            const recorder = recorderRef.current;

            if (recorder && recorder.state !== 'inactive') {
                // コンポーネント終了時はアップロード処理を開始しない
                recorder.onstop = null;
                recorder.ondataavailable = null;
                recorder.stop();
            }

            streamRef.current?.getTracks().forEach((track) => {
                track.stop();
            });
        };
    }, []);

    // お気に入り音声のメタデータを取得
    useEffect(() => {
        let cancelled = false;

        const loadMetadata = async () => {
            try {
                const metadata = await getVocabularyAudioMetadata(
                    vocabularyId,
                    audioType,
                );

                if (!cancelled) {
                    setIsFavorite(metadata.is_favorite);
                }
            } catch(err) {
                if (!cancelled) {
                    // 音声がまだ登録されていない場合
                    if (isAxiosError(err) && err.response?.status === 404) {
                        setIsFavorite(false);
                    } else {
                        setError('音声情報の取得に失敗しました。');
                    }
                }
            } finally {
                if (!cancelled) {
                    setFavoriteLoading(false);
                }
            }
        };

        void loadMetadata();

        return () => {
            cancelled = true;
        };
    }, [vocabularyId, audioType, audioVersion]);

    const startRecording = async () => {
        setError('');

        if (isFavorite) {
            setError('お気に入り登録中の音声は録音できません。');
            return;
        }

        if (!navigator.mediaDevices?.getUserMedia ||
            typeof MediaRecorder === 'undefined') {
                setError('このブラウザでは録音機能を利用できません。');
                return;
        }
        
        let stream: MediaStream | null = null;

        try {
            stream = await navigator.mediaDevices.getUserMedia({
                audio: true,
            });

            const supportedTypes = [
                'audio/webm;codecs=opus',
                'audio/ogg;codecs=opus',
                'audio/mp4',
            ];

            const mimeType = supportedTypes.find((type) =>
                MediaRecorder.isTypeSupported(type),
            );

            const recorder = mimeType
                ? new MediaRecorder(stream, { mimeType })
                : new MediaRecorder(stream);

            streamRef.current = stream;
            recorderRef.current = recorder;
            chunksRef.current = [];

            recorder.ondataavailable = (event) => {
                if (event.data.size > 0) {
                    chunksRef.current.push(event.data);
                }
            };

            recorder.onerror = () => {
                setError('録音中にエラーが発生しました。');
                setRecording(false);
                streamRef.current?.getTracks().forEach((track) => {
                    track.stop();
                });
                streamRef.current = null;
            };

            recorder.onstop = async () => {
                const blob = new Blob(chunksRef.current, {
                    type: recorder.mimeType || 'audio/webm',
                });

                streamRef.current?.getTracks().forEach((track) => {
                    track.stop();
                });
                streamRef.current = null;
                recorderRef.current = null;

                if (blob.size === 0) {
                    setError('音声が録音されませんでした。');
                    return;
                }

                if (blob.size > 10 * 1024 * 1024) {
                    setError('音声ファイルは 10MiB 以下にしてください。');
                    return;
                }

                setSaving(true);

                try {
                    await uploadVocabularyAudio(
                        vocabularyId,
                        audioType,
                        blob,
                    );

                    setFavoriteLoading(true);
                    setAudioVersion((version) => version + 1);
                } catch {
                    setError(
                        '音声の保存に失敗しました。通信やファイル形式を確認してください。',
                    );
                } finally {
                    setSaving(false);
                }
            };

            recorder.start();
            setRecording(true);
        } catch {
            stream?.getTracks().forEach((track) => {
                track.stop();
            });

            streamRef.current = null;
            setError(
                'マイクを使用できません。ブラウザのマイク許可設定を確認してください。',
            );
        }
    };

    const stopRecording = () => {
        const recorder = recorderRef.current;

        if (recorder && recorder.state !== 'inactive') {
            setRecording(false);
            recorder.stop();
        }
    };

    // お気に入りボタンの処理
    const toggleFavorite = async () => {
        setFavoriteLoading(true);
        setError('');

        const nextFavorite = !isFavorite;

        try {
            await setVocabularyAudioFavorite(
                vocabularyId,
                audioType,
                nextFavorite,
            );

            setIsFavorite(nextFavorite);
        } catch {
            setError('お気に入り状態の更新に失敗しました。');
        } finally {
            setFavoriteLoading(false);
        }
    };

    const audioUrl =
        `${getVocabularyAudioUrl(vocabularyId, audioType)}` +
        `&v=${audioVersion}`;
    
    return (
        <section
            style={{
                padding: '0',
                margin: '0',
                border: '1px solid #777',
                borderRadius: '8px',
                boxSizing: 'border-box',
                width: '100%',
            }}
        >
            <div
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    flexWrap: 'wrap',
                    gap: '8px',
                }}
            >
            <h4 style={{margin:'0', padding: '4px 0 0 10px'}}>{label}</h4>

            {/* お気に入りボタンの表示 */}
            <button
                type="button"
                onClick={toggleFavorite}
                disabled={favoriteLoading || recording || saving }
                style={{marginTop:'0', flexShrink: 0}}
            >
                {favoriteLoading
                     ? 'Loading,,,'
                     : isFavorite
                        ? '★ UnFAV'
                        : '☆ FAV'}
            </button>

            {!recording ? (
                <button style={{marginTop:'0', flexShrink: 0}}
                    type="button"
                    onClick={startRecording}
                    disabled={saving || isFavorite}
                >
                    {isFavorite ? 'LOCK' : saving ? 'Saving' : 'REC'}
                </button>
            ) : (
                <button style={{marginTop:'0', flexShrink: 0}}
                    type="button"
                    onClick={stopRecording}
                >
                    STOP
                </button>
            )}

            {recording && <p style={{flex: '1 1 100%', margin: 0}}>Recording</p>}
            {saving && <p style={{flex: '1 1 100%', margin: 0}}>Saving audio...</p>}
            </div>

            <div style={{ marginTop: '8px'}}>
                <audio
                    key={audioVersion}
                    controls
                    preload="none"
                    src={audioUrl}
                    style={{display: 'block', maxWidth: '100%'}}
                >
                    お使いのブラウザは音声再生に対応していません。
                </audio>
            </div>

            {error && (
                <p role="alert" style={{ color: "red"}}>
                    {error}
                </p>
            )}
        </section>
    );
};

export default VocabularyAudioControls;