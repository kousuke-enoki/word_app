import axiosInstance from '../../axiosConfig'

/**
 * 単語を登録するAPI
 * @param wordId 登録する単語のID
 * @param memo オプションのメモ
 */
export const saveMemo = async (wordId: number, memo = '') => {
  try {
    const response = await axiosInstance.post('/words/memo', {
      wordId,
      memo,
    })
    return response.data // 必要なら成功データを返す
  } catch (error) {
    console.error('Error saving memo:', error)
    throw error // 呼び出し元でエラーハンドリング
  }
}
