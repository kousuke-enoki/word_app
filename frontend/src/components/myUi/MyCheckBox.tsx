import { Checkbox } from '@headlessui/react'
import clsx from 'clsx'
import { Check } from 'lucide-react'

type Props = { checked: boolean; label: string; onChange: (v: boolean) => void }

export const MyCheckbox: React.FC<Props> = ({ checked, onChange, label }) => (
  <label className="flex cursor-pointer select-none items-center gap-2">
    <Checkbox
      checked={checked}
      onChange={onChange}
      className={clsx(
        'grid h-5 w-5 place-content-center rounded-md border transition focus:outline-none focus:ring-2 ring-[var(--button_bg)]',
        checked
          ? 'border-[var(--button_border)] bg-[var(--button_bg)] text-[var(--button)]'
          : 'border-[var(--border)] bg-[var(--container_bg)] text-[var(--fg)]',
      )}
    >
      <Check
        className={clsx(
          'h-4 w-4 transition',
          checked ? 'opacity-100' : 'opacity-0',
        )}
      />
    </Checkbox>
    <span className="text-[var(--fg)]">{label}</span>
  </label>
)
