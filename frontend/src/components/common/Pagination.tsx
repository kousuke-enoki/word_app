// src/components/common/Pagination.tsx
import React from 'react'

import { Button } from '@/components/ui/ui'

type Props = {
  /** 1-based current page */
  page: number
  /** total pages (>=1) */
  totalPages: number
  /** page change handler（1-basedで受け取ります） */
  onPageChange: (page: number) => void

  /** 現在のページサイズ（任意・指定時にドロップダウン表示） */
  pageSize?: number
  /** ページサイズ変更（任意） */
  onPageSizeChange?: (size: number) => void
  /** ページサイズ候補（任意, 例: [10,20,30,50]） */
  pageSizeOptions?: number[]

  className?: string
}

const Pagination: React.FC<Props> = ({
  page,
  totalPages,
  onPageChange,
  pageSize,
  onPageSizeChange,
  pageSizeOptions = [10, 20, 30, 50],
  className = '',
}) => {
  const clamp = (p: number) => Math.max(1, Math.min(totalPages, p))
  const go = (p: number) => onPageChange(clamp(p))

  return (
    <div className={`mt-4 flex flex-wrap items-center gap-2 ${className}`}>
      {/* ページサイズ切り替え（指定時のみ表示） */}
      {onPageSizeChange && typeof pageSize === 'number' && (
        <select
          className="rounded-xl border border-[var(--input_bd)] bg-[var(--select)] px-3 py-2 text-[var(--select_c)]"
          value={pageSize}
          onChange={(e) => onPageSizeChange(Number(e.target.value))}
        >
          {pageSizeOptions.map((n) => (
            <option key={n} value={n}>
              {n}
            </option>
          ))}
        </select>
      )}

      <Button onClick={() => go(1)} disabled={page === 1}>
        最初へ
      </Button>
      <Button onClick={() => go(page - 1)} disabled={page === 1}>
        前へ
      </Button>

      <span className="px-2 text-sm opacity-80">
        ページ {page} / {totalPages}
      </span>

      <Button onClick={() => go(page + 1)} disabled={page === totalPages}>
        次へ
      </Button>
      <Button onClick={() => go(totalPages)} disabled={page === totalPages}>
        最後へ
      </Button>
    </div>
  )
}

export default Pagination
