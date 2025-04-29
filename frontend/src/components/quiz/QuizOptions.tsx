import React from 'react'
 
type QuizOptionsProps = {
  onSelect: (testType: 'all' | 'registered' | 'custom' | 'start') => void
}

const QuizOptions: React.FC<QuizOptionsProps> = ({ onSelect }) => {
  return (
    <div>
      <h2>テストオプションを選択してください</h2>
      <button onClick={() => onSelect('all')}>すべての単語からテスト</button>
      <button onClick={() => onSelect('registered')}>登録単語からテスト</button>
      <button onClick={() => onSelect('custom')}>カスタムテスト</button>
    </div>
  )
}

export default QuizOptions
