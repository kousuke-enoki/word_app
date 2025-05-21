import { useState } from 'react';
import clsx from 'clsx';

type Props = {
  title: string;
  disabled: boolean;
  children: any;
};


export const MyCollapsible = ({title, disabled, children}: Props) => {
  const [open,setOpen]=useState(false);

  return (
    <div>
      <button onClick={()=>setOpen(!open)}
              className="flex w-full items-center justify-between rounded-md
                         bg-sky-100 px-3 py-2 text-sm font-medium text-sky-800
                         hover:bg-gray-200 transition
                         dark:bg-sky-800 text-white">
        <span>{title}</span>
        <svg className={clsx('h-4 w-4 transition-transform',
                             open ? 'rotate-180' : '')}
             viewBox="0 0 20 20" fill="currentColor">
          <path d="M6 8l4 4 4-4" />
        </svg>
      </button>
      {!disabled && open && <div className="mt-2 pl-3">{children}</div>}
    </div>
  );
};
