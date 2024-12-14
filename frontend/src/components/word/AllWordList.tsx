import React, { useState, useEffect } from 'react'
import axiosInstance from '../../axiosConfig'
import { useNavigate, useLocation } from 'react-router-dom'
import { Word, WordInfo, JapaneseMean } from '../../types/wordTypes'
import { registerWord } from '../../service/word/RegisterWord'
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
  const [isInitialized, setIsInitialized] = useState<boolean>(false)
  const [successMessage, setSuccessMessage] = useState<string>('')

  // APIからデータを取得する関数
  useEffect(() => {
    // location.state から値を復元（初回のみ実行）
    if (location.state) {
      setSearch(location.state.search || '')
      setSortBy(location.state.sortBy || 'name')
      setOrder(location.state.order || 'asc')
      setPage(location.state.page || 1)
      setLimit(location.state.limit || 10)
    }
    setIsInitialized(true) // 初期化完了をマーク
  }, [location.state])

  // 初回レンダリング時と依存する値が変わったときにデータを取得
  useEffect(() => {
    if (!isInitialized) return // 初期化が完了していなければ実行しない

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
  }, [search, sortBy, order, page, limit, isInitialized])

  // ページング処理
  const handlePageChange = (newPage: React.SetStateAction<number>) => {
    setPage(newPage)
  }

  // 詳細ページに遷移する関数
  const handleDetailClick = (id: number) => {
    navigate(`/words/${id}`, { state: { search, sortBy, order, page, limit } })
  }

  const handleRegister = async (word: Word) => {
    if (!word) return
    try {
      // API呼び出しから新しい登録状態と登録数を取得
      const updatedWord = await registerWord(word.id, !word.isRegistered)

      // 単語の登録状態を更新
      setWords((prevWords) =>
        prevWords.map((w) =>
          w.id === word.id
            ? {
                ...w,
                isRegistered: updatedWord.isRegistered,
                registrationCount: updatedWord.registrationCount,
              }
            : w,
        ),
      )
      const registeredWordName = updatedWord.name
      if (updatedWord.isRegistered) {
        setSuccessMessage(registeredWordName + ' を登録しました。')
      } else {
        setSuccessMessage(registeredWordName + ' を登録解除しました。')
      }
    } catch (error) {
      console.error('Error registering word:', error)
    }
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
      <select
        value={sortBy}
        onChange={(e) => {
          const newSortBy = e.target.value
          if (newSortBy === 'register' && sortBy !== 'register') {
            setPage(1)
          }
          setSortBy(newSortBy)
        }}
      >
        <option value="name">単語名</option>
        <option value="registrationCount">登録数</option>
        <option value="register">登録</option>
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
            <th>登録数</th>
            <th>登録</th>
            <th>詳細</th>
          </tr>
        </thead>
        <tbody>
          {words.map((word) => (
            <tr key={word.id}>
              <td className={`word-name`}>{word.name}</td>
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
              <td> {word.registrationCount} </td>
              <td>
                {successMessage && (
                  <div className="success-popup">{successMessage}</div>
                )}
                <div>
                  <button
                    className={`register-button ${word.isRegistered ? 'registered' : ''}`}
                    onClick={() => handleRegister(word)}
                  >
                    {word.isRegistered ? '解除' : '登録'}
                  </button>
                </div>
              </td>
              <td>
                <div>
                  <button onClick={() => handleDetailClick(word.id)}>
                    詳細
                  </button>
                </div>
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
