import { useEffect } from "react"

type Props = {
  open: boolean
  onClose: () => void
  children: React.ReactNode
  className?: string
}

function Modal({ open, onClose, children, className = "bg-white rounded-xl p-6 w-full max-w-md relative" }: Props) {
  if (!open) return null
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
        if (e.key === "Escape") onClose()
    }

    window.addEventListener("keydown", handler)

    return () => window.removeEventListener("keydown", handler)
    }, [])

    useEffect(() => {
  document.body.style.overflow = "hidden"
  return () => {
    document.body.style.overflow = "auto"
  }
}, [])

  return (
    <div 
      className="fixed inset-0 flex items-center justify-center bg-black/50 z-50 transition-opacity"
      onClick={onClose}
    >
      <div 
        className={className}
        onClick={(e) => e.stopPropagation()}
      >
        <button
          onClick={onClose}
          className="absolute top-2 right-2 text-gray-500"
        >
          ✕
        </button>

        {children}
      </div>
    </div>
  )
}

export default Modal