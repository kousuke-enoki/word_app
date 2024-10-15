import React, { useState, useEffect } from 'react';
import axiosInstance from '../axiosConfig';
import { resolveProjectReferencePath } from 'typescript';

// 単語の型定義

interface Word {
  id: number,
  name: string,
  edges: {
    WordInfos: WordInfo[],
  }
}

interface WordInfo {
  id: number,
  edges:{
    PartOfSpeech: PartOfSpeech
    JapaneseMeans: JapaneseMean[]
  }
}

interface PartOfSpeech {
  id: number,
  name: string,
}
interface JapaneseMean {
  id: number,
  name: string,
}

const AllWordList: React.FC = () => {
  const [words, setWords] = useState<Word[]>([]); // Word[] 型の単語リスト
  const [search, setSearch] = useState<string>('');
  const [sortBy, setSortBy] = useState<string>('name');
  const [order, setOrder] = useState<string>('asc');
  const [page, setPage] = useState<number>(1);
  const [totalPages, setTotalPages] = useState<number>(1);
  // const [limit, setLimit] = useState<number>(10);


// プロパティ名を変換する関数
const transformResponseData = (data: any): Word[] => {
  return data.map((word: any) => {
    const transformedWordInfos = word.edges.word_infos.map((wordInfo: any) => ({
      ...wordInfo,
      edges: {
        ...wordInfo.edges,
        PartOfSpeech: wordInfo.edges.part_of_speech,
        JapaneseMeans: wordInfo.edges.japanese_means,
      },
    }));

    return {
      ...word,
      edges: {
        WordInfos: transformedWordInfos,
      },
    };
  });
}

  // APIからデータを取得する関数
  const fetchWords = async () => {
    console.log(search)
    console.log(sortBy)
    console.log(order)
    console.log(page)
    try {
      const response = await axiosInstance.get('words/all_list', {
        params: {
          search,
          sortBy,
          order,
          page,
          // limit,
        },
      });
      console.log(response)
      console.log(response.data)
      const transformedData = transformResponseData(response.data.words)
      setWords(transformedData);
      setTotalPages(response.data.totalPages);
    } catch (error) {
      console.error('Failed to fetch words:', error);
    }
  };

  // 初回レンダリング時と依存する値が変わったときにデータを取得
  useEffect(() => {
    fetchWords();
  }, [search, sortBy, order, page]);

  // ページング処理
  const handlePageChange = (newPage: React.SetStateAction<number>) => {
    setPage(newPage);
  };

  return (
    <div>
      <h1>単語一覧</h1>

      {/* 検索フォーム */}
      <input
        type="text"
        placeholder="単語検索"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />

      {/* ソート選択 */}
      <select value={sortBy} onChange={(e) => setSortBy(e.target.value)}>
        <option value="name">単語名</option>
        {/* <option value="part_of_speech_id">品詞</option> */}
      </select>
      <button onClick={() => setOrder(order === 'asc' ? 'desc' : 'asc')}>
        {order === 'asc' ? '昇順' : '降順'}
      </button>

      {/* 単語のリスト表示 */}
      <table>
        <thead>
          <tr>
            <th>単語名</th>
            <th>日本語訳</th>
            <th>品詞</th>
            <th>編集</th>
            <th>詳細</th>
          </tr>
        </thead>
        <tbody>
          {words.map((word) => (
            <tr key={word.id}>
              <td>{word.name}</td>
              <td>{word.edges.WordInfos.map((info: WordInfo) => info.edges.JapaneseMeans.map(
                (JapaneseMean: JapaneseMean) => JapaneseMean.name).join(', ')).join(', ')}</td>
              <td>{word.edges.WordInfos.map((info: WordInfo) => info.edges.PartOfSpeech.name).join(', ')}</td>
              <td><button>編集</button></td>
              <td><button>詳細</button></td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* ページネーション */}
      <div>
        <button onClick={() => handlePageChange(1)} disabled={page === 1}>
          最初へ
        </button>
        <button onClick={() => handlePageChange(page - 1)} disabled={page === 1}>
          前へ
        </button>
        <span>ページ {page} / {totalPages}</span>
        <button onClick={() => handlePageChange(page + 1)} disabled={page === totalPages}>
          次へ
        </button>
        <button onClick={() => handlePageChange(totalPages)} disabled={page === totalPages}>
          最後へ
        </button>
      </div>
    </div>
  );
};

export default AllWordList;
