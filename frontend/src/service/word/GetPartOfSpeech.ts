export type PartOfSpeechOption = {
  id: number
  name: string
}
// 品詞のサンプルデータ
export const getPartOfSpeech: PartOfSpeechOption[] = [
  { id: 1, name: '名詞' },
  { id: 2, name: '代名詞' },
  { id: 3, name: '動詞' },
  { id: 4, name: '形容詞' },
  { id: 5, name: '副詞' },
  { id: 6, name: '助動詞' },
  { id: 7, name: '前置詞' },
  { id: 8, name: '冠詞' },
  { id: 9, name: '間投詞' },
  { id: 10, name: '接続詞' },
  { id: 11, name: '慣用句' },
  { id: 12, name: 'その他' },
]

 // フィルタリングされた品詞を取得する関数
 export const getPartsOfSpeechForQuiz = (): PartOfSpeechOption[] => {
  return getPartOfSpeech.filter((pos) => [1, 3, 4, 5].includes(pos.id))
}
