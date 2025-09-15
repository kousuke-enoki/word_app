import React from 'react'

type Props = {
  /** 見出しテキスト（または任意のノード） */
  title: React.ReactNode
  /** 右側に置きたい要素（任意：ボタン等） */
  right?: React.ReactNode
  /** 配置：中央寄せなら 'center'、左右に分けるなら 'between'  */
  align?: 'between' | 'center'
  className?: string
}

const PageTitle: React.FC<Props> = ({ title, right, align = 'between', className = '' }) => {
  const layout =
    align === 'center'
      ? 'justify-center'
      : 'justify-between'

  return (
    <div className={`mb-4 flex flex-wrap items-end gap-3 ${layout} ${className}`}>
      <h1 className="text-2xl font-bold text-[var(--h1_fg)]">{title}</h1>
      {align !== 'center' && right}
    </div>
  )
}

export default PageTitle
