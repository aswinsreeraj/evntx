type Props = {
  open: boolean
  onClose: () => void
  children: React.ReactNode
}

function Modal({ open, onClose, children }: Props) {
  if (!open) return null

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/50 z-50">
      <div className="bg-white rounded-xl p-6 w-full max-w-md relative">
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