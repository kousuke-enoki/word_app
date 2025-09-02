import React from 'react'


type CardProps<T extends React.ElementType = 'div'> = {
  as?: T
  className?: string
} & React.ComponentPropsWithoutRef<T>

export function Card<T extends React.ElementType = 'div'>({
  className = '',
  as,
  ...props
}: CardProps<T>) {
  const Tag = as || 'div'
  return (
    <Tag
      className={`rounded-2xl border border-[var(--border)] bg-[var(--container_bg)] text-[var(--container_c)] shadow-sm ${className}`}
      {...props}
    />
  )
}

type InputProps = React.InputHTMLAttributes<HTMLInputElement>
export const Input: React.FC<InputProps> = ({ className = '', ...props }) => (
  <input
    className={`w-full rounded-xl border border-[var(--input_bd)] bg-[var(--input)] px-4 py-2 text-[var(--input_c)] placeholder:opacity-60 outline-none focus:ring-2 ring-[var(--button_bg)] ${className}`}
    {...props}
  />
)


export const Badge: React.FC<React.HTMLAttributes<HTMLSpanElement>> = ({ className = '', ...props }) => (
  <span
    className={`inline-flex items-center gap-1 rounded-full border border-[var(--border)] bg-[var(--container_bg)] px-2 py-0.5 text-xs text-[var(--fg)] ${className}`}
    {...props}
  />
)


export const PageContainer: React.FC<React.HTMLAttributes<HTMLDivElement>> = ({ className = '', ...props }) => (
  <div className={`mx-auto w-full max-w-3xl px-4 py-10 ${className}`} {...props} />
)
