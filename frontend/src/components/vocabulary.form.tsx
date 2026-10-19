import { useState, type FormEvent } from "react";

import { createVocabulary } from "../services/vocabulary.service";

type Props = {
    onCreated: () => void;
};

const VocabularyForm = ({ onCreated }: Props) => {
    const [word, setWord] = useState('');

    const [meaning, setMeaning] = useState('');

    const [comment, setComment] = useState('');

    const [exampleSentence, setExampleSentence] = useState('');

    const [partOfSpeech, setPartOfSpeech] = useState('noun');

    const [error, setError]= useState('');

    const [submitting, setSubmitting] = useState(false);

    const submitHandler = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError('');

        if (!word.trim() || !meaning.trim()) {
            setError('単語と日本語訳を入力してください。');
            return;
        }

        try {
            setSubmitting(true);

            await createVocabulary({
                word: word.trim(),
                meaning: meaning.trim(),
                exampleSentence: exampleSentence.trim() || undefined,
                partOfSpeech,
                comment: comment.trim(),
            });

            setWord('');
            setMeaning('');
            setExampleSentence('');
            setComment('');

            await onCreated();
        } catch {
            setError('登録に失敗しました。入力内容やサーバーを確認してください。');
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <form onSubmit={submitHandler}>
            <h2>単語を登録</h2>

            <label>
                英単語
                <input
                    value={word}
                    onChange={(e) => setWord(e.target.value)}
                    required
                 />
            </label>

            <label>
                日本語訳
                <input
                    value={meaning}
                    onChange={(e) => setMeaning(e.target.value)}
                    required
                 />
            </label>

            <label>
                例文
                <textarea
                    value={exampleSentence}
                    onChange={(e) => setExampleSentence(e.target.value)}
                 />
            </label>

            <label>
                品詞
                <select
                    value={partOfSpeech}
                    onChange={(e) => setPartOfSpeech(e.target.value)}
                 >
                    <option value="noun">名詞</option>
                    <option value="verb">動詞</option>
                    <option value="adjective">形容詞</option>
                    <option value="adverb">副詞</option>
                    <option value="preposition">前置詞</option>
                    <option value="pronoun">代名詞</option>
                    <option value="conjunction">接続詞</option>
                    <option value="interjection">間投詞</option>
                </select>
            </label>

            <label>
                コメント
                <textarea
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                 />
            </label>

            {error && <p role="alert">{error}</p>}

             <button type="submit" disabled={submitting}>
                {submitting ? '登録中・・・' : '登録'}
             </button>
        </form>
    );
};

export default VocabularyForm;