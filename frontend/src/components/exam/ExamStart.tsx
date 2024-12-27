import React, { useEffect, useRef } from 'react'
import axiosInstance from '../../axiosConfig'

type ExamStartProps = {
  settings?: {
    questionCount: number
    targetWordTypes: string
    partsOfSpeeches: number[]
  }
}

const ExamStart: React.FC<ExamStartProps> = ({ settings }) => {
  const isFetchedRef = useRef(false) // フラグを追跡
  useEffect(() => {
    if (isFetchedRef.current) return // 2回目以降は実行しない
    isFetchedRef.current = true
    console.log(settings)
    // APIリクエストでテスト問題を取得
    const fetchExams = async () => {
      try {
        const response = await axiosInstance.post('/exams/new', settings)
        const data = await response.data
        console.log(data) // デバッグ用
      } catch (error) {
        console.log(error) // デバッグ用
        alert('テスト開始できませんでした。')
      }
    }

    fetchExams()
  }, [settings])

  return <h2>テストを開始します...</h2>
}

export default ExamStart
