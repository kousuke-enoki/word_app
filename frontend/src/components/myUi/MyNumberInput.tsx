import clsx from 'clsx'
import React from 'react'

type Props = {
  value: number
  min?: number
  max?: number
  onChange: (v: number) => void
}

export const MyNumberInput: React.FC<Props> = ({
  value,
  min,
  max,
  onChange,
}) => (
  <input
    type="number"
    min={min}
    max={max}
    value={value}
    onChange={(e) => onChange(Number(e.target.value))}
    className={clsx(
      'w-24 rounded-xl border px-3 py-2 text-right outline-none',
      'border-[var(--input_bd)] bg-[var(--input)] text-[var(--input_c)]',
      'focus:ring-2 ring-[var(--button_bg)]',
    )}
  />
)
