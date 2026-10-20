import { useCallback, useEffect, useState } from 'react';

import {
    fetchVocabularies
} from '../services/vocabulary.service';

import type { Vocabulary } from '../types/vocabulary';

import VocabularyForm from '../components/vocabulary.form';

const VocabularyPage = () => {
    const [vocabularies, setVocabularies] = useState<Vocabulary[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [keyword, setKeyword] = useState('');

    const loadVocabularies = useCallback(async () => {
        try {
            const data = await fetchVocabularies();
            setVocabularies(data);
            setError('');
        } catch {
            setError('単語一覧を取得できませんでした。');
        }
    }, []);

    useEffect(() => {
        let isCurrent = true;

        fetchVocabularies()
            .then((data) => {
                if (isCurrent) {
                    setVocabularies(data);
                }
            })
            .catch(() => {
                if (isCurrent) {
                    setError('単語一覧を取得できませんでした。');
                }
            })
            .finally(() => {
                if (isCurrent) {
                    setLoading(false);
                }
            });

        return () => {
            isCurrent = false;
        };
    }, []);

    const [sortOrder, setSortOrder] = useState('word-asc');

    const searchText = keyword.trim().toLowerCase();
    const filteredVocabularies = vocabularies.filter((vocabulary) => (
        vocabulary.word.toLowerCase().includes(searchText) ||
        vocabulary.meaning.toLowerCase().includes(searchText) ||
        (vocabulary.exampleSentence ?? '').toLowerCase().includes(searchText) ||
        vocabulary.comment.toLowerCase().includes(searchText)
    ));

    const sortedVocabularies = [...filteredVocabularies].sort((a, b) => {
        switch (sortOrder) {
            case 'word-desc':
                return b.word.localeCompare(a.word);

            case 'newest':
                return b.createdAt.localeCompare(a.createdAt);

            case 'oldest':
                return a.createdAt.localeCompare(b.createdAt);

            case 'word-asc':
            default:
                return a.word.localeCompare(b.word);
        }
    });

    return (
        <main style={{ border:"1px solid red"}}>
      <h1 style={{fontSize:"30px", backgroundColor:"green", borderRadius:"10px", width:"250px", margin:"20px auto", padding:"10px 0", color:"white", fontWeight:"bold"}}>英単語帳</h1>

      <VocabularyForm onCreated={loadVocabularies} />

      <h2>登録済みの単語</h2>

      <input
        type="search"
        placeholder="英単語・日本語訳・例文で検索"
        value={keyword}
        onChange={(e) => setKeyword(e.target.value)}
      />

      {loading && <p>読み込み中...</p>}
      {error && <p role="alert">{error}</p>}

      {!loading && !error && vocabularies.length === 0 && (
        <p>まだ単語が登録されていません。</p>
      )}

      {sortedVocabularies.map((vocabulary) => (
        <article key={vocabulary.id}>
          <h3>{vocabulary.word}</h3>
          <p>{vocabulary.meaning}</p>
          {vocabulary.exampleSentence && (
            <p>{vocabulary.exampleSentence}</p>
          )}
          <p>品詞：{vocabulary.partOfSpeech}</p>
          <p>コメント：{vocabulary.comment || 'なし'}</p>
        </article>
      ))}

      <select
        value={sortOrder}
        onChange={(e) => setSortOrder(e.target.value)}
        >
        <option value="word-asc">英単語 A → Z</option>
        <option value="word-desc">英単語 Z → A</option>
        <option value="newest">登録日が新しい順</option>
        <option value="oldest">登録日が古い順</option>
      </select>
    </main>
    );
};

export default VocabularyPage;