import { useCallback, useEffect, useState } from 'react';

import {
    fetchVocabularies,
    deleteVocabulary,
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
    // 品詞を絞り込む状態
    const [selectedPartOfSpeech, setSelectedPartOfSpeech] =
        useState('all');


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

    // 編集管理
    const [editingVocabulary, setEditingVocabulary] =
        useState<Vocabulary | null>(null);

    // 削除管理
    const [deletingVocabularyId, setDeletingVocabularyId] =
        useState<string | null>(null);
    const [actionError, setActionError] = useState('');

    const handleDelete = async (vocabulary: Vocabulary) => {
        if (!window.confirm(`「${vocabulary.word}」を削除しますか？`)) {
            return;
        }

        setDeletingVocabularyId(vocabulary.id);
        setActionError('');
        try {
            await deleteVocabulary(vocabulary.id);
            setVocabularies((current) =>
                current.filter((item) => item.id !== vocabulary.id)
            );
        } catch {
            setActionError('単語の削除に失敗しました。');
        } finally {
            setDeletingVocabularyId(null);
        }
    };
     
    // ソート条件の状態
    const [sortOrder, setSortOrder] = useState('word-asc');

    const searchText = keyword.trim().toLowerCase();
    const filteredVocabularies = vocabularies.filter((vocabulary) => {
        const matchesKeyword =
            vocabulary.word.toLowerCase().includes(searchText) ||
            vocabulary.meaning.toLowerCase().includes(searchText) ||
            (vocabulary.exampleSentence ?? '').toLowerCase().includes(searchText) ||
            vocabulary.comment.toLowerCase().includes(searchText);
        const matchesPartOfSpeech =
            selectedPartOfSpeech === 'all' ||
            vocabulary.partOfSpeech === selectedPartOfSpeech;

        return matchesKeyword && matchesPartOfSpeech;
    });

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

    <div style={{display:"flex", gap:"5px", flexWrap:"wrap", height:"200px" }}>
        {/* 単語登録フォーム */}
        <div style={{width:"700px", padding:"20px", margin:"20px 0",}}>
        <VocabularyForm onCreated={loadVocabularies}/>
        </div>
     {/* 検索・ソート・絞り込みフォーム */}
    <div style={{border:"1px solid blue", width:"350px", margin:"10px auto",}}>
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

      <label htmlFor='part-of-speech-filter'>
        品詞で絞り込む
      </label>

      <select
        id='part-of-speech-filter'
        value={selectedPartOfSpeech}
        onChange={(event) =>
            setSelectedPartOfSpeech(event.target.value)
        }
      >
        <option value="all">すべて</option>
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

      <p>
        表示件数：{sortedVocabularies.length}件
      </p>
      </div>
    </div>

        {/* 登録済み単語一覧 */}
      <h2>登録済みの単語</h2>
    <div style={{border:"1px solid red", display:"flex", gap:"5px", flexWrap:"wrap",}}>

      {loading && <p>読み込み中...</p>}
      {error && <p role="alert">{error}</p>}
      {actionError && <p role="alert">{actionError}</p>}

      {!loading && !error && vocabularies.length === 0 && (
        <p>まだ単語が登録されていません。</p>
      )}
      {!loading && !error && vocabularies.length > 0 && sortedVocabularies.length === 0 && (
        <p>条件に一致する単語がありません。</p>
      )}

      {sortedVocabularies.map((vocabulary) => (
        <div style={{border:"2px solid black", width:"300px", margin:"10px auto", borderRadius:"10px", background:"#a3d4d2"}}>
        <article
        key={vocabulary.id}>
          <h3 style={{backgroundColor:"white", fontWeight:"bold",}}>{vocabulary.word}</h3>
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
            <>
            <button style={{color:"green"}}
              type="button"
              onClick={() => setEditingVocabulary(vocabulary)}
            >
              編集
            </button>

            <button style={{color:"red"}}
                type='button'
                onClick={() => handleDelete(vocabulary)}
                disabled={deletingVocabularyId === vocabulary.id}
            >
                {deletingVocabularyId === vocabulary.id ? '削除中...' : '削除'}
            </button>
            </>
          )}
        </article>
        </div>
      ))}
    </div>

     
    </main>
    );
};

export default VocabularyPage;