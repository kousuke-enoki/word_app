import clsx from 'clsx'
import { Link } from 'react-router-dom'

import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { ResultQuestion } from '@/types/quiz'

type Props = {
  rows: ResultQuestion[]
  onToggleRegister: (row: ResultQuestion) => void
}

const ResultTable = ({ rows, onToggleRegister }: Props) => (
  <Card className="p-0 overflow-hidden">
    <table className="w-full border-separate border-spacing-0 text-sm">
      <thead>
        <tr className="bg-[var(--table_th)] text-[var(--table_th_c)]">
          <th className="px-3 py-2 text-left font-semibold rounded-l-lg">#</th>
          <th className="px-3 py-2 text-left font-semibold">単語</th>
          <th className="px-3 py-2 text-left font-semibold">正解</th>
          <th className="px-3 py-2 text-left font-semibold">選択</th>
          <th className="px-3 py-2 text-left font-semibold">登録</th>
          <th className="px-3 py-2 text-right font-semibold">回答回数</th>
          <th className="px-3 py-2 text-right font-semibold rounded-r-lg">
            正解回数
          </th>
        </tr>
      </thead>

      <tbody className="divide-y divide-[var(--border)]">
        {rows.map((row) => {
          const correct =
            row.choicesJpms.find((c) => c.japaneseMeanID === row.correctJpmID)
              ?.name ?? '-'
          const select =
            row.choicesJpms.find((c) => c.japaneseMeanID === row.answerJpmID)
              ?.name ?? '-'

          return (
            <tr
              key={row.questionNumber}
              className="
                transition-colors
                even:[&>td]:bg-[var(--table_tr_e)]
                hover:[&>td]:bg-[var(--table_row_hover)]
                active:[&>td]:bg-[var(--table_row_active)]
              "
            >
              <td className="px-3 py-2">{row.questionNumber}</td>

              <td className="px-3 py-2 underline-offset-2">
                <Link
                  to={`/words/${row.wordID}`}
                  className="text-[var(--primary)] underline"
                >
                  {row.wordName}
                </Link>
              </td>

              <td className="px-3 py-2">{correct}</td>

              <td
                className={clsx(
                  'px-3 py-2 rounded-md',
                  row.isCorrect
                    ? 'bg-emerald-50 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-200'
                    : 'bg-rose-50 text-rose-800 dark:bg-rose-900/30 dark:text-rose-200',
                )}
              >
                {select}
              </td>

              <td className="px-3 py-2">
                {row.registeredWord.isRegistered ? (
                  <Button
                    onClick={() => onToggleRegister(row)}
                    className="px-2 py-1 text-xs !bg-rose-600 !border-rose-600 !text-white hover:opacity-90"
                  >
                    解除
                  </Button>
                ) : (
                  <Button
                    onClick={() => onToggleRegister(row)}
                    className="px-2 py-1 text-xs !bg-emerald-600 !border-emerald-600 !text-white hover:opacity-90"
                  >
                    登録
                  </Button>
                )}
              </td>

              <td className="px-3 py-2 text-right">
                {row.registeredWord.quizCount}
              </td>
              <td className="px-3 py-2 text-right">
                {row.registeredWord.correctCount}
              </td>
            </tr>
          )
        })}
      </tbody>
    </table>
  </Card>
)

export default ResultTable
