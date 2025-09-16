import {
  autoUpdate,
  flip,
  offset,
  shift,
  size,
  useFloating,
} from '@floating-ui/react'
import { Listbox, Portal, Transition } from '@headlessui/react'
import { Check, ChevronsUpDown } from 'lucide-react'
import { Fragment } from 'react'

type Option<T> = { value: T; label: string }
type Props<T> = {
  options: Option<T>[]
  value: T
  onChange: (v: T) => void
  className?: string
}

export function MySelect<T extends string | number>({
  options,
  value,
  onChange,
  className = 'w-52',
}: Props<T>) {
  const selected = options.find((o) => o.value === value) ?? options[0]

  // ▼ トリガーとパネルの位置合わせ
  const { refs, floatingStyles } = useFloating({
    placement: 'bottom-start',
    strategy: 'fixed', // スクロールでも安定
    middleware: [
      offset(6),
      flip(),
      shift({ padding: 8 }),
      // パネル幅をボタン幅に合わせる
      size({
        apply({ rects, elements, availableWidth }) {
          Object.assign(elements.floating.style, {
            minWidth: `${rects.reference.width}px`,
            maxWidth: `${availableWidth}px`,
          })
        },
      }),
    ],
    whileElementsMounted: autoUpdate,
  })

  return (
    <Listbox value={value} onChange={onChange}>
      <div className={`relative mt-1 ${className}`}>
        {/* Trigger */}
        <Listbox.Button
          ref={refs.setReference}
          className="
            listbox-btn relative w-full cursor-pointer rounded-md
            !bg-[var(--select)] !text-[var(--select_c)]
            py-2 pl-3 pr-10 text-left shadow-sm
            !border !border-[var(--select_bd)]
            ring-1 ring-inset !ring-[var(--select_bd)]
            focus:outline-none focus:!ring-2 focus:!ring-[var(--button_bg)]
            transition-colors
          "
          style={{
            backgroundColor: 'var(--select, #fff)',
            color: 'var(--select_c, #111)',
            borderColor: 'var(--select_bd, #cfd4dc)',
          }}
        >
          <span className="block truncate">{selected?.label}</span>
          <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2 opacity-80">
            <ChevronsUpDown size={18} />
          </span>
        </Listbox.Button>

        <Transition
          as={Fragment}
          enter="transition ease-out duration-100"
          enterFrom="opacity-0 scale-95"
          enterTo="opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="opacity-100 scale-100"
          leaveTo="opacity-0 scale-95"
        >
          <Portal>
            {/* Dropdown */}
            <Listbox.Options
              ref={refs.setFloating}
              style={{
                ...floatingStyles, // ← 位置を適用
                backgroundColor: 'var(--select, #fff)',
                color: 'var(--select_c, #111)',
                borderColor: 'var(--select_bd, #cfd4dc)',
                zIndex: 9999,
              }}
              className="
                listbox-panel max-h-60 overflow-auto rounded-md
                !bg-[var(--select)] !text-[var(--select_c)]
                py-1 shadow-lg !border
                focus:outline-none
              "
            >
              {options.map((opt) => (
                <Listbox.Option
                  key={opt.value}
                  value={opt.value}
                  className={({ active }) => `
                    relative cursor-pointer select-none py-2 pl-10 pr-4
                    ${
                      active
                        ? // ← ホバーしても色を変えない（背景も同色で固定）
                          '!bg-[var(--select)] !text-[var(--select_c)]'
                        : '!text-[var(--select_c)]'
                    }
                  `}
                >
                  {({ selected }) => (
                    <>
                      <span
                        className={`block truncate ${selected ? 'font-medium' : 'font-normal'}`}
                      >
                        {opt.label}
                      </span>
                      {selected && (
                        <span className="absolute inset-y-0 left-0 flex items-center pl-3 !text-[var(--button)]">
                          <Check size={18} />
                        </span>
                      )}
                    </>
                  )}
                </Listbox.Option>
              ))}
            </Listbox.Options>
          </Portal>
        </Transition>
      </div>
    </Listbox>
  )
}
