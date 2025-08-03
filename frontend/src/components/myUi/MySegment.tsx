import clsx from 'clsx';

type Props = {
  value: number;
  onChange: (v: number) => void;
  targets: { value: number; label: string }[];
};

export const MySegment = ({ value, onChange, targets }: Props) => (
  <div role="radiogroup">
    {targets.map((o, i) => {
      const active = value === o.value;
      return (
        <button
          key={o.value}
          type="button"
          role="radio"
          aria-checked={active}
          onClick={() => onChange(o.value)}
          className={clsx(
            "px-4 py-1 text-sm font-semibold transition-colors duration-150 focus:outline-none",
            // 丸みは先頭/末尾だけ強めに
            i === 0 && "rounded-l-md",
            i === targets.length - 1 && "rounded-r-md",
            active
              ? "bg-blue-600 text-white shadow-inner"
              : "bg-blue-200 text-black hover:bg-blue-300 dark:bg-blue-800 dark:text-blue-100 dark:hover:bg-blue-700"
          )}
        >
          {o.label}
        </button>
      );
    })}
  </div>
);
