export interface ResultSummary {
  quizNumber: number
  createdAt: string            // ISO-8601
  isRegisteredWords: number    // 0:全単語 1:登録のみ 2:未登録のみ
  isIdioms: number             // 0:すべて 1:慣用句のみ 2:除外
  isSpecialCharacters: number  // 0:すべて 1:特殊のみ 2:除外
  choicesPosIds: number[]      // 品詞 ID 配列
  totalQuestionsCount: number
  correctCount: number
  resultCorrectRate: number    // 0-100 %
}
