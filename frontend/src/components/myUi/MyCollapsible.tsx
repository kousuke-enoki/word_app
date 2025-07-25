import clsx from 'clsx';
import { useId, useState } from 'react';

type Props = {
  title: string;
  disabled?: boolean;
  children: React.ReactNode;
  defaultOpen?: boolean;
  removeCollapsedGap?: boolean; // true なら閉じた時に下の帯を完全に消す
};

export const MyCollapsible: React.FC<Props> = ({
  title,
  disabled = false,
  children,
  defaultOpen = false,
  removeCollapsedGap = false,
}) => {
  const [open, setOpen] = useState(defaultOpen);
  const id = useId();

  return (
    <section
      className={clsx(
        'collapsible',
        removeCollapsedGap && !open && 'border-b-0' /* 任意 */,
      )}
      data-open={open}
      aria-disabled={disabled}
    >
      <button
        type="button"
        className="collapsible__header"
        aria-expanded={open}
        aria-controls={id}
        disabled={disabled}
        onClick={() => !disabled && setOpen(o => !o)}
      >
        <span>{title}</span>
        <svg className="collapsible__icon" viewBox="0 0 10 10" fill="currentColor">
          <path d="M6 8l4 4 4-4" />
        </svg>
      </button>

      {/* 非表示にして空白を完全になくす場合 */}
      {removeCollapsedGap ? (
        open && (
          <div id={id}>
            <div className="collapsible__inner">{children}</div>
          </div>
        )
      ) : (
        <div id={id} className="collapsible__content">
          <div className="collapsible__inner">{children}</div>
        </div>
      )}
    </section>
  );
};
