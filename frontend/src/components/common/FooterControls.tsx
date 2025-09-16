// src/components/common/FooterControls.tsx
import React from 'react'

import PageBottomNav, {
  PageBottomNavProps,
} from '@/components/common/PageBottomNav'
import Pagination, { PaginationProps } from '@/components/common/Pagination'
import { Card } from '@/components/ui/card'

type Props = {
  pagination: PaginationProps
  nav: PageBottomNavProps
  className?: string
  /** テーブルの直下に薄く“浮く”感じにしたい時 */
  sticky?: boolean
}

/**
 * テーブル下に置く操作フッター
 * - md以上: 左にページネーション / 右に下部ナビ（縦仕切り線）
 * - モバイル: 2段に積む（間に水平の仕切り線）
 */
const FooterControls: React.FC<Props> = ({
  pagination,
  nav,
  className = '',
  sticky,
}) => (
  <Card
    className={[
      'mt-6 p-4',
      'backdrop-blur supports-[backdrop-filter]:bg-[color:var(--container_bg)/.92]',
      sticky ? 'sticky bottom-4 z-20' : '',
      className,
    ].join(' ')}
  >
    <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      {/* 左：ページネーション（コンパクト） */}
      <div className="min-w-0">
        <Pagination {...pagination} compact className="m-0" />
      </div>

      {/* 仕切り線（画面幅で出し分け） */}
      <div className="h-px bg-[var(--border)] md:hidden" />
      <div className="hidden h-6 w-px self-stretch bg-[var(--border)] md:block" />

      {/* 右：下部ナビ（横並び & コンパクト） */}
      <div className="min-w-0 md:justify-end">
        <PageBottomNav {...nav} inline compact />
      </div>
    </div>
  </Card>
)

export default FooterControls
