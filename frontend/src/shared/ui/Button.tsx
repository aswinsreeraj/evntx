type Props = {
  children: React.ReactNode
  onClick?: () => void
  className?: string
}

function Button({ children, onClick, className }: Props) {
  return (
    <button
      onClick={onClick}
      className={`px-4 py-2 rounded-lg bg-black text-white hover:bg-gray-800 transition ${className}`}
    >
      {children}
    </button>
  )
}

export default Button