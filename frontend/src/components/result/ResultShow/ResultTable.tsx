import clsx from 'clsx'
import { Link } from 'react-router-dom'

import Pagination, { PaginationProps } from '@/components/common/Pagination'
import { RegisterToggle } from '@/components/common/RegisterToggle'
import { Card } from '@/components/ui/card'
import { ResultQuestion } from '@/types/quiz'

type PagerInTable = PaginationProps & { align?: 'left' | 'center' | 'right' }

type Props = {
  rows: ResultQuestion[]
  onToggleRegister: (row: ResultQuestion) => void
  /** テーブル下フッターに出すページネーション（省略可） */
  pager?: PagerInTable
}

const ResultTable = ({ rows, onToggleRegister, pager }: Props) => (
  <Card className="p-0">
    {/* ▼ 横スクロール・モバイル */}
    <div className="overflow-x-auto touch-pan-x overscroll-x-contain">
      <table className="w-full table-fixed border-separate border-spacing-0 text-sm min-w-[720px] sm:min-w-[860px]">
        <thead>
          <tr className="bg-[var(--table_th)] text-[var(--table_th_c)]">
            {/* ▼ 各列に幅（レスポンシブ）を与える */}
            <th className="w-10 sm:w-12 px-2 py-2 text-left font-semibold rounded-l-lg">
              #
            </th>
            <th className="w-40 sm:w-56 px-3 py-2 text-left font-semibold">
              単語
            </th>
            <th className="w-36 sm:w-44 px-3 py-2 text-left font-semibold">
              正解
            </th>
            <th className="w-36 sm:w-44 px-3 py-2 text-left font-semibold">
              選択
            </th>
            <th className="w-28 sm:w-40 px-3 py-2 text-left font-semibold">
              登録
            </th>
            <th className="w-20 px-2 py-2 text-right font-semibold">
              回答回数
            </th>
            <th className="w-20 px-2 py-2 text-right font-semibold rounded-r-lg">
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
                {/* ID は小さく・中央寄せ */}
                <td className="px-2 py-2 text-center">{row.questionNumber}</td>

                {/* 単語は truncate（table-fixed なので max-w-0 を併用） */}
                <td className="px-3 py-2 max-w-0 truncate">
                  <Link
                    to={`/words/${row.wordID}`}
                    className="underline underline-offset-2 text-blue-600 hover:opacity-80 dark:text-sky-300"
                    title={row.wordName}
                  >
                    {row.wordName}
                  </Link>
                </td>

                {/* 正解／選択も truncate */}
                <td className="px-3 py-2 max-w-0 truncate" title={correct}>
                  {correct}
                </td>

                <td
                  className={clsx(
                    'px-3 py-2 max-w-0 truncate rounded-md',
                    row.isCorrect
                      ? 'bg-emerald-50 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-200'
                      : 'bg-rose-50 text-rose-800 dark:bg-rose-900/30 dark:text-rose-200',
                  )}
                  title={select}
                >
                  {select}
                </td>

                {/* 登録：固定幅のトグル（モバイルで少し小さめ） */}
                <td className="px-3 py-2">
                  <RegisterToggle
                    isRegistered={row.registeredWord.isRegistered}
                    onToggle={() => onToggleRegister(row)}
                    variant="compact"
                    widthClass="w-24 sm:w-28" // ← 小画面で幅を詰める
                  />
                </td>

                <td className="px-2 py-2 text-right">
                  {row.registeredWord.quizCount}
                </td>
                <td className="px-2 py-2 text-right">
                  {row.registeredWord.correctCount}
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>

    {/* ページネーション */}
    {pager && (
      <div className="border-t border-[var(--border)]">
        {/* ← text-xxx ではなく flex + justify-xxx で位置決め */}
        <div
          className={clsx(
            'border-t border-[var(--border)] px-2 py-0', // ← 余白をゼロに
            'leading-none', // ← 行高を 1 に固定
          )}
        >
          <Pagination {...pager} compact className="!mt-0 !mb-0" />
        </div>
      </div>
    )}
  </Card>
)

export default ResultTable
