// src/components/common/Modal.tsx
import React from 'react'

type Props = {
  open: boolean
  onClose: () => void
  title?: string
  children?: React.ReactNode
  widthClass?: string // 任意
}

const Modal: React.FC<Props> = ({
  open,
  onClose,
  title,
  children,
  widthClass,
}) => {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <div
        className={`relative z-10 ${widthClass ?? 'w-[min(720px,92vw)]'} rounded-2xl bg-[var(--container_bg)] p-4 shadow-xl`}
      >
        {title && <h3 className="mb-3 text-lg font-semibold">{title}</h3>}
        {children}
      </div>
    </div>
  )
}

export default Modal
