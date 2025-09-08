import React from 'react'

import { Badge } from '@/components/ui/card'

type Props = { correct: number; total: number; rate: number }

const ResultHeader: React.FC<Props> = ({ correct, total, rate }) => (
  <div className="mb-4 flex flex-wrap items-end gap-3">
    <h1 className="text-2xl font-bold text-[var(--h1_fg)]">クイズ結果</h1>
    <div className="flex items-center gap-2">
      <Badge>
        正解 {correct}/{total}
      </Badge>
      <Badge>{rate.toFixed(1)}%</Badge>
    </div>
  </div>
)

export default ResultHeader
