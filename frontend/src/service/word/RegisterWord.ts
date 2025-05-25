import axiosInstance from '../../axiosConfig'

/**
 * 単語を登録するAPI
 * @param wordId 登録する単語のID
 * @param IsRegistered 登録するか、解除するか
 */
export const registerWord = async (wordId: number, IsRegistered: boolean) => {
  try {
    const response = await axiosInstance.post('/words/register', {
      wordId,
      IsRegistered,
    })
    return response.data // 必要なら成功データを返す
  } catch (error) {
    console.error('Error registering word:', error)
    throw error // 呼び出し元でエラーハンドリング
  }
}
