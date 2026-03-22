import { useEffect } from "react"

type Props = {
  open: boolean
  onClose: () => void
  children: React.ReactNode
  className?: string
}

function Modal({ open, onClose, children, className = "bg-white rounded-xl p-6 w-full max-w-md relative" }: Props) {
  useEffect(() => {
    if (!open) return

    const handler = (e: KeyboardEvent) => {
        if (e.key === "Escape") onClose()
    }

    window.addEventListener("keydown", handler)

    return () => window.removeEventListener("keydown", handler)
    }, [open, onClose])

  if (!open) return null

  return (
    <div
      className="fixed inset-0 flex items-center justify-center bg-black/50 z-50 transition-opacity"
      onClick={onClose}
    >
      <div
        className={className}
        onClick={(e) => e.stopPropagation()}
      >

        {children}
      </div>
    </div>
  )
}

export default Modal
