import { Switch } from '@headlessui/react';
import clsx from 'clsx';

type Props = {
  checked: boolean;
  onChange: (v: boolean) => void;
  id: string;
};


export const MySwitch = ({ checked, onChange, id }: Props) => (
  <Switch
    id={id}
    checked={checked}
    onChange={onChange}
    className={clsx(
      'appearance-none p-0 border-0 relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
      checked ? 'bg-blue-600' : 'bg-gray-300'
    )}
  >
    <span
      className={clsx(
        'inline-block h-5 w-5 transform rounded-full bg-white shadow transition',
        checked ? 'translate-x-5' : 'translate-x-1'
      )}
    />
  </Switch>
);
