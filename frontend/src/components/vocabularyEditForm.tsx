import { useState } from "react";
import type { SubmitEvent } from "react";

import type { Vocabulary } from "../types/vocabulary";
import { updateVocabulary } from "../services/vocabulary.service";

type VocabularyEditFormProps = {
    vocabulary: Vocabulary;
    onUpdated: () => void;
    onCancel: () => void;
};

const VocabularyEditForm = ({
    vocabulary,
    onUpdated,
    onCancel,
}: VocabularyEditFormProps) => {
    const [word, setWord] = useState(vocabulary.word);
    const [meaning, setMeaning] = useState(vocabulary.meaning);
    const [exampleSentence, setExampleSentence] = useState(
        vocabulary.exampleSentence ?? ''
    );
    const [partOfSpeech, setPartOfSpeech] = useState(
        vocabulary.partOfSpeech
    );
    const [comment, setComment] = useState(vocabulary.comment);

    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (event: SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();

        try {
            setSaving(true);
            setError('');

            await updateVocabulary(vocabulary.id, {
                word,
                meaning,
                exampleSentence: exampleSentence || undefined,
                partOfSpeech,
                comment,
            });

            onUpdated();
        } catch {
            setError('単語の更新に失敗しました。');
        } finally {
            setSaving(false);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
        <h3>単語を編集</h3>
        
        <div>
            <label htmlFor="edit-word">英単語</label>
            <input
                id="edit-word"
                value={word}
                onChange={(event) => setWord(event.target.value)}
                required
             />
        </div>

        <div>
            <label htmlFor="edit-meaning">日本語訳</label>
            <input
                id="edit-meaning"
                value={meaning}
                onChange={(event) => setMeaning(event.target.value)}
                required
             />
        </div>

        <div>
            <label htmlFor="edit-example">例文</label>
            <input
                id="edit-example"
                value={exampleSentence}
                onChange={(event) => setExampleSentence(event.target.value)}
             />
        </div>

        <div>
            <label htmlFor="edit-part-of-speech">品詞</label>
            <select
                id="edit-part-of-speech"
                value={partOfSpeech}
                onChange={(event) => setPartOfSpeech(event.target.value)}
                required
             >
                <option value="noun">名詞</option>
                <option value="verb">動詞</option>
                <option value="adjective">形容詞</option>
                <option value="adverb">副詞</option>
                <option value="preposition">前置詞</option>
                <option value="pronoun">代名詞</option>
                <option value="conjunction">接続詞</option>
                <option value="interjection">間投詞</option>
                <option value="enclitic">小辞</option>
            </select>
        </div>

        <div>
            <label htmlFor="edit-comment">コメント</label>
            <input
                id="edit-comment"
                value={comment}
                onChange={(event) => setComment(event.target.value)}
             />
        </div>

        {error && <p role="alert">{error}</p>}

        <button type="submit" disabled={saving}>
            {saving ? '保存中...' : '保存'}
        </button>
        <button type="button" onClick={onCancel} disabled={saving}>
            キャンセル
        </button>
        </form>
    );
};

export default VocabularyEditForm;