import { Listbox, Transition } from '@headlessui/react';
import { Check, ChevronsUpDown } from 'lucide-react';
import { Fragment } from 'react';

type Option<T> = { value: T; label: string };
type Props<T> = { options: Option<T>[]; value: T; onChange: (v: T) => void };

export function MySelect<T extends string | number>({ options, value, onChange }: Props<T>) {
  const selected = options.find(o => o.value === value)!;

  return (
    <Listbox value={value} onChange={onChange}>
      <div className="relative mt-1 w-52">
        <Listbox.Button
          className="
            listbox-btn relative w-full cursor-pointer rounded-md
            bg-white py-2 pl-3 pr-10 text-left shadow-sm
            ring-1 ring-inset ring-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500
            text-gray-900 dark:bg-gray-700 dark:text-white
          "
        >
          <span className="block truncate">{selected.label}</span>
          <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2 text-gray-500 dark:text-gray-300">
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
          <Listbox.Options
            className="
              absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md
              bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5
              dark:bg-gray-700
            "
          >
            {options.map(opt => (
              <Listbox.Option
                key={opt.value}
                value={opt.value}
                className={({ active }) =>
                  `relative cursor-pointer select-none py-2 pl-10 pr-4 ${
                    active
                      ? 'bg-blue-600 text-white'
                      : 'text-gray-900 dark:text-gray-100'
                  }`
                }
              >
                {({ selected }) => (
                  <>
                    <span className={`block truncate ${selected ? 'font-medium' : 'font-normal'}`}>
                      {opt.label}
                    </span>
                    {selected && (
                      <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-white dark:text-blue-200">
                        <Check size={18} />
                      </span>
                    )}
                  </>
                )}
              </Listbox.Option>
            ))}
          </Listbox.Options>
        </Transition>
      </div>
    </Listbox>
  );
}
