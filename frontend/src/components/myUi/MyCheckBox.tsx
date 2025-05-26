import { Checkbox } from '@headlessui/react';
import clsx from 'clsx';

type Props = {
  checked: boolean;
  label: string;
  onChange: (v: boolean) => void;
};

export const MyCheckbox: React.FC<Props> = ({ checked, onChange, label }) => (
  <label className="flex items-center gap-2 cursor-pointer select-none">
    <Checkbox
      checked={checked}
      onChange={onChange}
      className={clsx(
        'h-5 w-5 shrink-0 rounded border',
        checked
          ? 'border-blue-600 bg-blue-600 text-white'
          : 'border-gray-300 bg-white',
        'focus:outline-none focus:ring-2 focus:ring-blue-500'
      )}
    >
      {checked && <span className="block text-center text-xs font-bold">✓</span>}
    </Checkbox>
    <span>{label}</span>
  </label>
);
