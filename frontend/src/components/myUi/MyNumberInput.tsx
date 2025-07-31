import clsx from 'clsx';
import React from 'react';

type Props = {
  value: number;
  min?: number;
  max?: number;
  onChange: (v: number) => void;
};

export const MyNumberInput: React.FC<Props> = ({ value, min, max, onChange }) => (
  <input
    type="number"
    min={min}
    max={max}
    value={value}
    onChange={(e) => onChange(Number(e.target.value))}
    className={clsx(
      'w-24 rounded-md border border-gray-300 bg-white p-1 text-right',
      'focus:border-blue-500 focus:ring-blue-500 text-gray-900',
      'dark:bg-gray-700 dark:border-gray-600 dark:text-white'
    )}
  />
);
