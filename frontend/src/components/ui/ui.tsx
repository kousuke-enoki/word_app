import React from 'react'

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'ghost' | 'outline'
  full?: boolean
}

export const Button: React.FC<ButtonProps> = ({
  className = '',
  variant = 'primary',
  full,
  ...props
}) => {
  const base =
    'inline-flex items-center justify-center gap-2 rounded-xl px-4 py-2 text-sm font-medium transition focus:outline-none focus:ring-2 ring-[var(--button_bg)]'
  const variants: Record<string, string> = {
    primary:
      '!text-[var(--button)] !bg-[var(--button_bg)] !border !border-[var(--button_border)] hover:opacity-90 shadow-sm',
    ghost:
      'bg-transparent text-[var(--fg)] border border-transparent hover:bg-[var(--container_bg)]',
    outline:
      '!bg-[var(--btn-subtle-bg)] text-[var(--fg)] !border !border-[var(--btn-subtle-bd)] hover:opacity-95',
  }
  return (
    <button
      className={`${base} ${variants[variant]} ${full ? 'w-full' : ''} ${className}`}
      {...props}
    />
  )
}
