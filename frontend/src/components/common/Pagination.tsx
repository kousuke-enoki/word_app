import React from 'react'

import { Button } from '@/components/ui/ui'

const PAGE_SIZES = [10, 30, 50, 100] as const
type PageSize = (typeof PAGE_SIZES)[number]

type Props = {
  sizes?: readonly PageSize[]
  size: number
  page: number
  total: number
  onSize: (s: PageSize) => void
  onPrev: () => void
  onNext: () => void
}

const Pagination = ({
  sizes = PAGE_SIZES,
  size,
  page,
  total,
  onSize,
  onPrev,
  onNext,
}: Props) => (
  <div className="flex items-center justify-between">
    {/* ▼ ページサイズ（セグメント） */}
    <div className="inline-flex rounded-xl bg-[var(--btn-subtle-bd)] p-0.5 gap-0.5">
      {sizes.map((s) => (
        <Button
          key={s}
          variant={s === size ? 'primary' : 'ghost'}
          className="rounded-lg px-3 py-1.5 min-w-12 text-center"
          onClick={() => onSize(s)}
        >
          {s}
        </Button>
      ))}
    </div>

    {/* ▼ Prev / Next */}
    <div className="rs-pager space-x-2">
      <Button
        variant="outline"
        disabled={page === 0}
        onClick={onPrev}
        className="px-3"
      >
        Prev
      </Button>
      <Button
        variant={
          page + 1 >= Math.max(1, Math.ceil(total / size))
            ? 'outline'
            : 'primary'
        }
        disabled={page + 1 >= Math.max(1, Math.ceil(total / size))}
        onClick={onNext}
        className="px-3"
      >
        Next
      </Button>
    </div>
  </div>
)

export default Pagination
