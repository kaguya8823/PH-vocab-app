import { useCallback, useEffect, useState } from 'react';

import {
    fetchVocabularies,
} from '../services/vocabulary.service';

import type { Vocabulary } from '../types/vocabulary';

import VocabularyForm from '../components/vocabulary.form';
import VocabularyEditForm from '../components/vocabularyEditForm';

const VocabularyPage = () => {
    // 単語一覧の状態
    const [vocabularies, setVocabularies] = useState<Vocabulary[]>([]);
    // 読み込み・エラーの状態
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    // 検索キーワードの状態
    const [keyword, setKeyword] = useState('');

    // APIから単語一覧を取得
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

    // 編集・削除管理
    const [editingVocabulary, setEditingVocabulary] =
        useState<Vocabulary | null>(null);
     
    // ソート条件の状態
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

        {/* 単語登録フォーム */}
        <VocabularyForm onCreated={loadVocabularies} />

     {/* 検索・ソートフォーム */}   
    <div style={{border:"1px solid blue", width:"350px", margin:"10px auto", }}>
      <input
        type="search"
        placeholder="キーワードで検索"
        style={{width:"330px", padding:"10px", margin:"10px 0", fontSize:"15px"}}
        value={keyword}
        onChange={(e) => setKeyword(e.target.value)}
      />
      <select
        value={sortOrder}
        style={{padding:"10px", margin:"10px 0", }}
        onChange={(e) => setSortOrder(e.target.value)
        }
        >
        <option value="word-asc">英単語 A → Z</option>
        <option value="word-desc">英単語 Z → A</option>
        <option value="newest">登録日が新しい順</option>
        <option value="oldest">登録日が古い順</option>
      </select>
      </div>

        {/* 登録済み単語一覧 */}
      <h2>登録済みの単語</h2>

      {loading && <p>読み込み中...</p>}
      {error && <p role="alert">{error}</p>}

      {!loading && !error && vocabularies.length === 0 && (
        <p>まだ単語が登録されていません。</p>
      )}

      {sortedVocabularies.map((vocabulary) => (
        <article
        key={vocabulary.id}>
          <h3>{vocabulary.word}</h3>
          <p>{vocabulary.meaning}</p>
          {vocabulary.exampleSentence && (
            <p>{vocabulary.exampleSentence}</p>
          )}
          <p>品詞：{vocabulary.partOfSpeech}</p>
          <p>コメント：{vocabulary.comment || 'なし'}</p>
          {editingVocabulary?.id === vocabulary.id ? (
            <VocabularyEditForm
              vocabulary={editingVocabulary}
              onUpdated={async () => {
                await loadVocabularies();
                setEditingVocabulary(null);
              }}
              onCancel={() => setEditingVocabulary(null)}
            />
          ) : (
            <button
              type="button"
              onClick={() => setEditingVocabulary(vocabulary)}
            >
              編集
            </button>
          )}
        </article>
      ))}

     
    </main>
    );
};

export default VocabularyPage;