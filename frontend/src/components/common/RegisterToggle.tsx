import clsx from 'clsx'
import React from 'react'

import { Button } from '@/components/ui/ui'

type Props = {
  isRegistered: boolean
  onToggle: () => void
  className?: string
  /** compact: ピル1つでトグル / split: 状態＋ボタンを分離 */
  variant?: 'compact' | 'split'
  /** テーブル内で安定するようデフォ固定サイズ */
  widthClass?: string // 例: "w-28"
}

export const RegisterToggle: React.FC<Props> = ({
  isRegistered,
  onToggle,
  className,
  variant = 'compact', // ← splitにすると状態＋ボタン表示
  widthClass = 'w-28', // ←固定幅。全行同じ幅で行ぶれ防止
}) => {
  if (variant === 'split') {
    return (
      <div
        className={clsx(
          'inline-grid grid-cols-2 gap-2 items-center',
          widthClass,
          className,
        )}
      >
        <span
          className={clsx(
            'inline-flex items-center justify-center rounded-full px-2.5 py-1 text-xs font-medium whitespace-nowrap',
            isRegistered
              ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-200'
              : 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-200',
          )}
        >
          {isRegistered ? '登録済み' : '未登録'}
        </span>
        {isRegistered ? (
          <Button
            variant="outline"
            onClick={onToggle}
            className="h-7 px-2 py-1 text-xs whitespace-nowrap"
            aria-label="登録を解除"
          >
            登録解除
          </Button>
        ) : (
          <Button
            onClick={onToggle}
            className="h-7 px-2 py-1 text-xs whitespace-nowrap"
            aria-label="登録する"
          >
            登録する
          </Button>
        )}
      </div>
    )
  }

  // compact: ピル1つで状態＋操作（クリックでトグル）
  return (
    <button
      type="button"
      onClick={onToggle}
      className={clsx(
        'inline-flex items-center justify-center rounded-full font-medium transition whitespace-nowrap',
        'h-7 sm:h-8 text-[11px] sm:text-xs',
        widthClass,
        isRegistered
          ? 'bg-blue-100 text-blue-800 border border-blue-200 dark:bg-blue-900/30 dark:text-blue-200 dark:border-blue-800'
          : '!bg-[var(--button_bg)] !text-[var(--button)] border border-[var(--button_border)] hover:opacity-90',
      )}
      aria-pressed={isRegistered}
      aria-label={isRegistered ? '登録解除' : '登録する'}
      title={isRegistered ? 'クリックで登録解除' : 'クリックで登録する'}
    >
      {isRegistered ? '登録済み' : '登録する'}
    </button>
  )
}
