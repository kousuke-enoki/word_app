import React, { useState, useEffect } from 'react';
import axiosInstance from '../../axiosConfig';
import { QuizQuestion, ChoiceJpm, AnswerRouteRes } from '../../types/quiz';
import '../../styles/components/quiz/QuizQuestionView.css'

interface Props {
  /** 表示・回答対象の問題 */
  question: QuizQuestion;
  /** 回答後コールバック。nextQuestion または finish result を渡す */
  onAnswered: (res: AnswerRouteRes) => void;
}

const QuizQuestionView: React.FC<Props> = ({ question, onAnswered }) => {
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [posting, setPosting] = useState(false);

  /* question が変わったら state をリセット */
  useEffect(() => {
    setSelectedId(null);
    setPosting(false);
  }, [question.quizID, question.questionNumber]);

  const handleSubmit = async () => {
    if (selectedId == null) return;
    setPosting(true);
    try {
      const payload = {
        quizID: question.quizID,
        answerJpmID: selectedId,
        questionNumber: question.questionNumber,
      };
      const res = await axiosInstance.post<AnswerRouteRes>(
        `/quizzes/answers/${question.quizID}`,
        payload,
      );
      onAnswered(res.data);
    } catch (e) {
      console.error(e);
      alert('回答送信に失敗しました');
      setPosting(false);
    }  finally {
      setPosting(false);
    }
  };

  /** 2×2 グリッド用に行列化 */
  const gridChoices = [...question.choicesJpms]
    .concat(new Array(4).fill({ japaneseMeanID: -1, name: '' }))
    .slice(0, 4);
  while (gridChoices.length < 4) gridChoices.push({ japaneseMeanID: -1, name: '' });

  const renderChoice = (c: ChoiceJpm) => {
    const isSelected = selectedId === c.japaneseMeanID;
    return (
      <button
        key={c.japaneseMeanID}
        disabled={c.japaneseMeanID === -1}
        onClick={() => setSelectedId(c.japaneseMeanID)}
        className={`question-button ${isSelected ? 'selected' : ''}`}
      >
        {c.name}
      </button>
    );
  };

  return (
    <div className="mx-auto max-w-xl space-y-8">
      {/* 上段 */}
      <div className="flex justify-between">
        <span className="text-sm font-semibold">
          Q{question.questionNumber}
        </span>
      </div>

      {/* 単語 */}
      <h1 className="text-center text-4xl font-bold">{question.wordName}</h1>

      {/* 選択肢 2×2 */}
      <div className="grid grid-cols-2 gap-4">
        {gridChoices.map(renderChoice)}
      </div>

      {/* OK ボタン */}
      <div className="text-center">
        <button
          onClick={handleSubmit}
          disabled={selectedId == null || posting}
          className="rounded bg-blue-600 px-6 py-2 font-medium text-white
                     disabled:bg-gray-400"
        >
          OK
        </button>
      </div>
    </div>
  );
};

export default QuizQuestionView;
