// src/components/common/Pagination.tsx
import React from 'react'

import { Button } from '@/components/ui/ui'

export type PaginationProps = {
  page: number
  totalPages: number
  onPageChange: (page: number) => void
  pageSize?: number
  onPageSizeChange?: (size: number) => void
  pageSizeOptions?: number[]
  className?: string
  compact?: boolean
}

const Pagination: React.FC<PaginationProps> = ({
  page,
  totalPages,
  onPageChange,
  pageSize,
  onPageSizeChange,
  pageSizeOptions = [10, 20, 30, 50],
  className = '',
  compact = false,
}) => {
  const clamp = (p: number) => Math.max(1, Math.min(totalPages, p))
  const go = (p: number) => onPageChange(clamp(p))

  const btnSize = compact ? 'px-3 py-1 text-sm rounded-lg' : ''
  const selSize = compact
    ? 'px-2 py-1 text-sm rounded-lg'
    : 'px-3 py-2 rounded-xl'
  const gapCls = compact ? 'gap-1' : 'gap-2'

  const SelectNode =
    onPageSizeChange != null && typeof pageSize === 'number' ? (
      <select
        data-testid="pagination-page-size" // 統合テストでページサイズ選択要素を特定するため
        className={`border border-[var(--input_bd)] bg-[var(--select)] text-[var(--select_c)] ${selSize}`}
        value={pageSize}
        onChange={(e) => onPageSizeChange(Number(e.target.value))}
      >
        {pageSizeOptions.map((n) => (
          <option key={n} value={n}>
            {n}
          </option>
        ))}
      </select>
    ) : null

  const Core = (
    <div className={`flex flex-wrap items-center justify-center ${gapCls}`}>
      <Button className={btnSize} onClick={() => go(1)} disabled={page === 1}>
        最初へ
      </Button>
      <Button
        className={btnSize}
        onClick={() => go(page - 1)}
        disabled={page === 1}
      >
        前へ
      </Button>
      <span
        className={`px-2 ${compact ? 'text-xs leading-none' : 'text-sm'} opacity-80`}
      >
        ページ {page} / {totalPages}
      </span>
      <Button
        className={btnSize}
        onClick={() => go(page + 1)}
        disabled={page === totalPages}
      >
        次へ
      </Button>
      <Button
        className={btnSize}
        onClick={() => go(totalPages)}
        disabled={page === totalPages}
      >
        最後へ
      </Button>
    </div>
  )

  // ラッパーは最小余白＆中央寄せ
  return (
    <div
      className={`mx-auto w-full flex px-2 py-1 justify-center ${className}`}
    >
      {/* xs: 1列(=2行) / sm+: 3カラム */}
      <div className="grid grid-cols-1 sm:grid-cols-[auto_1fr_auto] items-center gap-1 leading-none">
        {/* セレクト：xsでは下段センター、sm+では左端 */}
        <div className="order-2 mt-1 justify-self-center sm:order-1 sm:mt-0 sm:justify-self-end">
          {SelectNode}
        </div>

        {/* Core：常に中央 */}
        <div className="order-1 justify-self-center sm:order-2">{Core}</div>

        {/* 右のダミーselect：sm+のみ表示してセンタリングを厳密化 */}
        <div className="hidden sm:block order-3 justify-self-start">
          {SelectNode ? (
            <div
              className="opacity-0 h-0 overflow-hidden pointer-events-none"
              aria-hidden
            >
              {SelectNode}
            </div>
          ) : null}
        </div>
      </div>
    </div>
  )
}

export default Pagination
