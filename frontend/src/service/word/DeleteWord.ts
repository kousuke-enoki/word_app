import axiosInstance from '../../axiosConfig'

/**
 * 単語を削除するAPI
 * @param wordId 削除する単語のID
 */

export const deleteWord = async (wordId: number): Promise<void> => {
  try {
    await axiosInstance.delete(`/words/${wordId}`)
  } catch (error) {
    console.error('Error deleting word:', error)
    throw error
  }
}
