// src/components/common/Modal.tsx
import React, { useEffect } from 'react'

type Props = {
  open: boolean
  onClose: () => void
  title?: string
  children?: React.ReactNode
  widthClass?: string
  overlayClassName?: string
  panelClassName?: string
}

const Modal: React.FC<Props> = ({
  open,
  onClose,
  title,
  children,
  widthClass,
  overlayClassName,
  panelClassName,
}) => {
  useEffect(() => {
    if (!open) return
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      document.body.style.overflow = prev
    }
  }, [open])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-[1000] flex items-center justify-center">
      <div
        className={overlayClassName ?? 'absolute inset-0 bg-black/70'}
        onClick={onClose}
      />
      <div
        role="dialog"
        aria-modal="true"
        className={
          // ← 背景/枠/影は CSS 変数で切り替え。ライトは白、ダークは黒系に自動で。
          `relative z-10 ${widthClass ?? 'w-[min(720px,92vw)]'} rounded-2xl
           bg-[var(--modal-bg)]
           border border-[var(--modal-border)]
           shadow-[var(--modal-shadow)]
           p-4` + (panelClassName ? ` ${panelClassName}` : '')
        }
      >
        {title && <h3 className="mb-3 text-lg font-semibold">{title}</h3>}
        {children}
      </div>
    </div>
  )
}

export default Modal
