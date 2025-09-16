import clsx from 'clsx'
import { useId, useState } from 'react'

type Props = {
  title: string
  disabled?: boolean
  children: React.ReactNode
  defaultOpen?: boolean
  removeCollapsedGap?: boolean
}

export const MyCollapsible: React.FC<Props> = ({
  title,
  disabled = false,
  children,
  defaultOpen = false,
  removeCollapsedGap = false,
}) => {
  const [open, setOpen] = useState(defaultOpen)
  const id = useId()

  return (
    <section
      className={clsx(
        'space-y-3',
        disabled && 'opacity-60 pointer-events-none',
      )}
    >
      <button
        type="button"
        aria-expanded={open}
        aria-controls={id}
        onClick={() => !disabled && setOpen((o) => !o)}
        className={clsx(
          'flex w-full items-center justify-between rounded-xl border px-3 py-2 text-left',
          'transition focus:outline-none focus:ring-2 ring-[var(--button_bg)]',
          'bg-[var(--btn-subtle-bg)] border-[var(--btn-subtle-bd)] text-[var(--fg)]',
        )}
      >
        <span className="font-medium">{title}</span>
        <svg
          viewBox="0 0 20 20"
          className={clsx(
            'h-4 w-4 transform transition-transform',
            open && 'rotate-180',
          )}
          fill="currentColor"
        >
          <path d="M5.5 8l4.5 4 4.5-4" />
        </svg>
      </button>

      {open ? (
        <div
          id={id}
          className="rounded-xl border border-[var(--border)] bg-[var(--container_bg)] p-4"
        >
          {children}
        </div>
      ) : (
        !removeCollapsedGap && <div className="h-2" />
      )}
    </section>
  )
}
