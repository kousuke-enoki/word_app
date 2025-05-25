/** MySegment.tsx : 2-option トグル */

type Props = {
  value: number;
  onChange: (v: number) => void;
  targets: {
    value: number;
    label: string;
  }[]
};

export const MySegment = ({value, onChange, targets }: Props) => (
  <div className="inline-flex rounded-lg bg-gray-200 p-1 dark:bg-gray-700">
    {targets.map(o => (
      <button
        key={o.value}
        onClick={()=>onChange(o.value)}
        className={`px-4 py-1 rounded-md transition
          ${value===o.value ? 'bg-blue-600 text-white' : 'text-gray-700'}`}
      >
        {o.label}
      </button>
    ))}
  </div>
);
