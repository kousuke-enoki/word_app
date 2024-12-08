import React, { useState, useEffect } from 'react'
import axiosInstance from '../../axiosConfig'
import { useNavigate, useLocation } from 'react-router-dom'
import { Word, WordInfo, JapaneseMean } from '../../types/wordTypes'
import '../../styles/components/word/AllWordList.css'

const AllWordList: React.FC = () => {
  const [words, setWords] = useState<Word[]>([])
  const [search, setSearch] = useState<string>('')
  const [sortBy, setSortBy] = useState<string>('name')
  const [order, setOrder] = useState<string>('asc')
  const location = useLocation()
  const [page, setPage] = useState<number>(location.state?.page || 1)
  const [totalPages, setTotalPages] = useState<number>(1)
  const [limit, setLimit] = useState<number>(10)
  const navigate = useNavigate()

  // 初回レンダリング時と依存する値が変わったときにデータを取得
  useEffect(() => {
    // APIからデータを取得する関数
    const fetchWords = async () => {
      try {
        const response = await axiosInstance.get('words/all_list', {
          params: {
            search,
            sortBy,
            order,
            page,
            limit,
          },
        })
        setWords(response.data.words)
        setTotalPages(response.data.totalPages)
      } catch (error) {
        console.error('Failed to fetch words:', error)
      }
    }
    fetchWords()
  }, [search, sortBy, order, page, limit])

  // ページング処理
  const handlePageChange = (newPage: React.SetStateAction<number>) => {
    setPage(newPage)
  }

  // 詳細ページに遷移する関数
  const handleDetailClick = (id: number) => {
    navigate(`/words/${id}`, { state: { page } })
  }

  return (
    <div className="wordList-container">
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
              <td>
                {word.wordInfos
                  .map((info: WordInfo) =>
                    info.japaneseMeans
                      .map((japaneseMean: JapaneseMean) => japaneseMean.name)
                      .join(', '),
                  )
                  .join(', ')}
              </td>
              <td>
                {word.wordInfos
                  .map((info: WordInfo) => info.partOfSpeech.name)
                  .join(', ')}
              </td>
              <td>
                <button>編集</button>
              </td>
              <td>
                <button onClick={() => handleDetailClick(word.id)}>詳細</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className="pagination-container">
        <select
          className="select-limit"
          value={limit}
          onChange={(e) => setLimit(Number(e.target.value))}
        >
          <option value="10">10</option>
          <option value="20">20</option>
          <option value="30">30</option>
          <option value="50">50</option>
        </select>

        <button onClick={() => handlePageChange(1)} disabled={page === 1}>
          最初へ
        </button>
        <button
          onClick={() => handlePageChange(page - 1)}
          disabled={page === 1}
        >
          前へ
        </button>
        <span>
          ページ {page} / {totalPages}
        </span>
        <button
          onClick={() => handlePageChange(page + 1)}
          disabled={page === totalPages}
        >
          次へ
        </button>
        <button
          onClick={() => handlePageChange(totalPages)}
          disabled={page === totalPages}
        >
          最後へ
        </button>
      </div>
    </div>
  )
}

export default AllWordList
