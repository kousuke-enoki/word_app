// MySegment.tsx
import clsx from 'clsx'

type Props = {
  value: number
  onChange: (v: number) => void
  targets: { value: number; label: string }[]
}

export const MySegment: React.FC<Props> = ({ value, onChange, targets }) => (
  <div
    role="radiogroup"
    className="
      inline-flex rounded-xl p-0.5 gap-0.5
      bg-[var(--btn-subtle-bd)]              /* 仕切り色を親に持たせる */
    "
  >
    {targets.map((o) => {
      const active = value === o.value
      return (
        <button
          key={o.value}
          type="button"
          role="radio"
          aria-checked={active}
          onClick={() => onChange(o.value)}
          className={clsx(
            'px-4 py-1.5 text-sm font-medium transition focus:outline-none focus:ring-2 ring-[var(--button_bg)]',
            'rounded-lg', // 子の角丸
            active
              ? '!bg-[var(--button_bg)] !text-[var(--button)]'
              : 'bg-[var(--btn-subtle-bg)] text-[var(--fg)] hover:opacity-95',
          )}
        >
          {o.label}
        </button>
      )
    })}
  </div>
)
