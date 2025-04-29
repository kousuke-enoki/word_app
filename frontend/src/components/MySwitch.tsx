// src/components/MySwitch.tsx
import * as Switch from '@radix-ui/react-switch';
import React from 'react';

type MySwitchProps = {
  /** 親が保持する boolean 状態 */
  checked: boolean;
  /** 親から渡されるハンドラ */
  onCheckedChange: (checked: boolean) => void;
};

export const MySwitch: React.FC<MySwitchProps> = ({
  checked,
  onCheckedChange,
}) => (
  <Switch.Root
    checked={checked}
    onCheckedChange={onCheckedChange}
    className="w-10 h-6 bg-gray-300 rounded-full relative data-[state=checked]:bg-blue-500 transition-colors"
  >
    <Switch.Thumb
      className="block w-4 h-4 bg-white rounded-full shadow-md translate-x-1 data-[state=checked]:translate-x-5 transition-transform"
    />
  </Switch.Root>
);
