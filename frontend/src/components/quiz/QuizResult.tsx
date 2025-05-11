// src/components/quiz/QuizResult.tsx
import React, { useMemo, useState } from 'react';
import axiosInstance from '../../axiosConfig';
import {
  ChoiceJpm,
  QuizQuestion,
  ResultRes,
  ResultQuestion,
} from '../../types/quiz';
import { registerWord } from '../../service/word/RegisterWord'
import { useNavigate } from 'react-router-dom';

interface Props {
  result: ResultRes;
}

const pageSizes = [10, 20, 30] as const;

const QuizResult: React.FC<Props> = ({ result }) => {
  const navigate = useNavigate();
  const [pageSize, setPageSize] = useState<(typeof pageSizes)[number]>(10);
  const [page, setPage] = useState(0);
  const [successMessage, setSuccessMessage] = useState<string>('')
  const [rowsData, setRowsData] = useState<ResultQuestion[]>(
    result.resultQuestions,
  );
  const [flash, setFlash] = useState('');
  const rows = useMemo(() => {
    const start = page * pageSize;
    return rowsData.slice(start, start + pageSize);
  }, [rowsData, page, pageSize]);

  console.log(result)

  /* ---------- 登録/解除ボタン ---------- */
  const handleRegister = async (idx: number) => {
    const q = rows[idx];
    try {
      const u = await registerWord(q.wordID, !q.registeredWord.isRegistered);

      // 行を更新
      setRowsData((prev) =>
        prev.map((row, i) =>
          i === idx
            ? {
                ...row,
                registeredWord: {
                  ...row.registeredWord,
                  isRegistered: u.isRegistered,
                  quizCount: u.quizCount,
                  correctCount: u.correctCount,
                },
              }
            : row,
        ),
      );

      setFlash(
        `${q.wordName} を${
          u.isRegistered ? '登録しました' : '登録解除しました'
        }`,
      );
      setTimeout(() => setFlash(''), 2500);
    } catch (e) {
      console.error(e);
      setFlash('通信に失敗しました');
      setTimeout(() => setFlash(''), 2500);
    }
  };

  const updateAttention = async (word: string, level: number) => {
    try {
      await axiosInstance.post('/registered-words/attention', { word, level });
    } catch (e) {
      console.error(e);
    }
  };

  /* ---------- UI ---------- */
  return (
    <div className="mx-auto max-w-4xl space-y-8 px-4 py-6">
      {/* サマリ */}
      <h1 className="text-2xl font-bold text-center">
        {result.correctCount}/{result.totalQuestionsCount} 正解
        <span className="ml-2 text-lg text-gray-500">
          ({(result.resultCorrectRate * 100).toFixed(1)}%)
        </span>
      </h1>

      {/* 設定表示 */}
      <div className="rounded-md border p-4 text-sm">
        <p>登録単語: {result.resultSetting.isRegisteredWords}</p>
        <p>慣用句: {result.resultSetting.isIdioms}</p>
        <p>特殊単語: {result.resultSetting.isSpecialCharacters}</p>
      </div>

      {/* 明細表 */}
      <table className="w-full border-collapse text-sm">
        <thead>
          <tr className="bg-gray-100">
            <th className="border p-2">#</th>
            <th className="border p-2">単語</th>
            <th className="border p-2">選択</th>
            <th className="border p-2">正誤</th>
            <th className="border p-2">登録</th>
            <th className="border p-2">回答回数</th>
            <th className="border p-2">正解回数</th>
            <th className="border p-2">注意Lv.</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((q: ResultQuestion, i: number) => (
            <tr key={q.questionNumber}>
              <td className="border p-1 text-center">{q.questionNumber}</td>
              <td className="border p-1">{q.wordName}</td>
              <td className="border p-1">
                {
                  q.choicesJpms.find(
                    (c) => c.japaneseMeanID === q.answerJpmId,
                  )?.name
                }
              </td>
              <td className="border p-1 text-center">
                {q.isCorrect ? '〇' : '×'}
              </td>
              {/* <td className="border p-1 text-center">
                <button
                  onClick={() => toggleRegister(q.wordName)}
                  className="rounded bg-green-600 px-2 py-1 text-white hover:bg-green-700"
                >
                  登録
                </button>
              </td> */}
              {/* <td>
                {successMessage && (
                  <div className="success-popup">{successMessage}</div>
                )}
                <div>
                  <button
                    className={`register-button ${q.registeredWord.isRegistered ? 'registered' : ''}`}
                    onClick={() => handleRegister(q)}
                  >
                    {q.registeredWord.isRegistered ? '解除' : '登録'}
                  </button>
                </div>
              </td> */}
              <td className="border p-1 text-center">
                {successMessage && (
                  <div className="success-popup">{successMessage}</div>
                )}
                <div>
                  <button
                    onClick={() => handleRegister(i + page * pageSize)}
                    className={`rounded px-2 py-1 text-white ${
                      q.registeredWord.isRegistered
                        ? 'bg-red-600 hover:bg-red-700'
                        : 'bg-green-600 hover:bg-green-700'
                    }`}
                  >
                    {q.registeredWord.isRegistered ? '解除' : '登録'}
                  </button>
                </div>
              </td>
              <td className="border p-1 text-center">{q.registeredWord.quizCount ?? 0}</td>
              <td className="border p-1 text-center">{q.registeredWord.correctCount ?? 0}</td>
              <td className="border p-1">
                <select
                  defaultValue={q.registeredWord.attentionLevel ?? 1}
                  onChange={(e) =>
                    updateAttention(q.wordName, Number(e.target.value))
                  }
                  className="rounded border px-1"
                >
                  {[1, 2, 3, 4, 5].map((l) => (
                    <option key={l}>{l}</option>
                  ))}
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* ページネーション */}
      <div className="flex items-center justify-between">
        <div className="space-x-2">
          {pageSizes.map((s) => (
            <button
              key={s}
              className={`rounded px-3 py-1 ${
                s === pageSize ? 'bg-blue-600 text-white' : 'bg-gray-200'
              }`}
              onClick={() => {
                setPageSize(s);
                setPage(0);
              }}
            >
              {s}
            </button>
          ))}
        </div>
        <div className="space-x-2">
          <button
            disabled={page === 0}
            onClick={() => setPage((p) => p - 1)}
            className="rounded bg-gray-200 px-3 py-1 disabled:opacity-50"
          >
            Prev
          </button>
          <button
            disabled={
              (page + 1) * pageSize >= result.resultQuestions.length
            }
            onClick={() => setPage((p) => p + 1)}
            className="rounded bg-gray-200 px-3 py-1 disabled:opacity-50"
          >
            Next
          </button>
        </div>
      </div>

      {/* ナビゲーション */}
      <div className="mt-6 flex justify-center gap-4">
        <button
          onClick={() => navigate('/quiz')}
          className="rounded bg-blue-500 px-4 py-2 text-white"
        >
          クイズメニュー
        </button>
        <button
          onClick={() => navigate('/result')}
          className="rounded bg-blue-500 px-4 py-2 text-white"
        >
          成績一覧
        </button>
        <button
          onClick={() => navigate('/')}
          className="rounded bg-gray-500 px-4 py-2 text-white"
        >
          ホーム
        </button>
      </div>
    </div>
  );
};

export default QuizResult;
