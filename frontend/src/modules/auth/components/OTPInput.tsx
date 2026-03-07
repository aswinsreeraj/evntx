import { useRef } from "react"

type Props = {
  value: string
  onChange: (value: string) => void
  length?: number
}

export default function OTPInput({ value, onChange, length = 6 }: Props) {
  const inputs = useRef<(HTMLInputElement | null)[]>([])

  const handleChange = (index: number, digit: string) => {
    if (!/^\d?$/.test(digit)) return

    const newOtp = value.split("")
    newOtp[index] = digit

    const otpString = newOtp.join("")
    onChange(otpString)

    if (digit && index < length - 1) {
      inputs.current[index + 1]?.focus()
    }
  }

  const handleKeyDown = (index: number, e: React.KeyboardEvent) => {
    if (e.key === "Backspace" && !value[index] && index > 0) {
      inputs.current[index - 1]?.focus()
    }
  }

  const handlePaste = (e: React.ClipboardEvent) => {
    const paste = e.clipboardData.getData("text").slice(0, length)
    if (!/^\d+$/.test(paste)) return

    onChange(paste)

    paste.split("").forEach((digit, index) => {
      if (inputs.current[index]) {
        inputs.current[index]!.value = digit
      }
    })

    inputs.current[paste.length - 1]?.focus()
  }

  return (
    <div className="flex gap-2 w-full justify-between" onPaste={handlePaste}>
      {Array.from({ length }).map((_, index) => (
        <input
          key={index}
          ref={(el) => (inputs.current[index] = el)}
          maxLength={1}
          placeholder="0"
          className="w-12 h-12 text-center border border-gray-300 rounded-xl text-lg font-medium text-gray-700 placeholder-gray-300 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-colors"
          style={{ 
             borderColor: value[index] ? '#a855f7' : '' // active/filled state purple outline based on design
          }}
          value={value[index] || ""}
          onChange={(e) => handleChange(index, e.target.value)}
          onKeyDown={(e) => handleKeyDown(index, e)}
        />
      ))}
    </div>
  )
}