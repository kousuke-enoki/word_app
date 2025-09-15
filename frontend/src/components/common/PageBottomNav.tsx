// src/components/common/PageBottomNav.tsx
import React from 'react'
import { useNavigate } from 'react-router-dom'

import { Button } from '@/components/ui/ui'

export type PageBottomNavProps = {
  actions?: {
    label: string
    to: string
    variant?: 'primary' | 'outline' | 'ghost'
  }[]
  back?: {
    label: string
    to: string
    variant?: 'primary' | 'outline' | 'ghost'
  }
  showHome?: boolean
  className?: string
  /** 背景/枠なしの横並び */
  inline?: boolean
  /** 小さめUI */
  compact?: boolean
}

const PageBottomNav: React.FC<PageBottomNavProps> = ({
  actions = [],
  back,
  showHome,
  className = '',
  inline,
  compact,
}) => {
  const nav = useNavigate()
  const btnSize = compact ? 'px-3 py-1.5 text-sm rounded-lg' : ''

  const Wrap: React.FC<React.HTMLAttributes<HTMLDivElement>> = ({
    children,
  }) =>
    inline ? (
      <div
        className={`flex flex-wrap items-center justify-center gap-2 ${className}`}
      >
        {children}
      </div>
    ) : (
      <div
        className={`rounded-2xl border border-[var(--border)] bg-[var(--container_bg)] p-4 flex flex-wrap items-center justify-center gap-3 ${className}`}
      >
        {children}
      </div>
    )

  return (
    <Wrap>
      {actions.map((a) => (
        <Button
          key={a.to}
          variant={a.variant ?? 'primary'}
          onClick={() => nav(a.to)}
          className={btnSize}
        >
          {a.label}
        </Button>
      ))}
      {back && (
        <Button
          variant={back.variant ?? 'outline'}
          onClick={() => nav(back.to)}
          className={btnSize}
        >
          {back.label}
        </Button>
      )}
      {showHome && (
        <Button variant="ghost" onClick={() => nav('/')} className={btnSize}>
          ホームに戻る
        </Button>
      )}
    </Wrap>
  )
}

export default PageBottomNav
