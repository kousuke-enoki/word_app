// src/components/ui/ui.tsx
import React from 'react'

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'ghost' | 'outline'
  full?: boolean
}

export const Button: React.FC<ButtonProps> = ({
  className = '',
  variant = 'primary',
  full,
  disabled,
  ...props
}) => {
  const base = [
    'inline-flex items-center justify-center gap-2 rounded-xl px-4 py-2 text-sm font-medium',
    'transition focus:outline-none focus:ring-2 ring-[var(--button_bg)]',
    // disabled の共通ふるまい
    'disabled:cursor-not-allowed disabled:pointer-events-none',
    'whitespace-nowrap leading-tight',
  ].join(' ')

  // disabled のときは hover クラス自体を出さない
  const hovPrimary = disabled ? '' : 'hover:opacity-90'
  const hovGhost = disabled ? '' : 'hover:bg-[var(--container_bg)]'
  const hovOutline = disabled ? '' : 'hover:opacity-95'

  const variants: Record<string, string> = {
    primary: [
      '!text-[var(--button)] !bg-[var(--button_bg)] !border !border-[var(--button_border)]',
      hovPrimary,
      'shadow-sm',
      // disabled の固定色
      'disabled:!bg-[var(--button_bg_disabled)]',
      'disabled:!text-[var(--button_disabled)]',
      'disabled:!border-[var(--button_border_disabled)]',
    ].join(' '),

    ghost: [
      'bg-transparent text-[var(--fg)] border border-transparent',
      hovGhost,
      'disabled:!text-[var(--button_disabled)]',
      'disabled:border-transparent',
      'disabled:bg-transparent',
    ].join(' '),

    outline: [
      '!bg-[var(--btn-subtle-bg)] text-[var(--fg)] !border !border-[var(--btn-subtle-bd)]',
      hovOutline,
      'disabled:!bg-[var(--pagination_bg_disabled,var(--btn-subtle-bg))]',
      'disabled:!text-[var(--button_disabled)]',
      'disabled:!border-[var(--pagination_bd,var(--btn-subtle-bd))]',
    ].join(' '),
  }

  return (
    <button
      disabled={disabled}
      className={`${base} ${variants[variant]} ${full ? 'w-full' : ''} ${className}`}
      {...props}
    />
  )
}
