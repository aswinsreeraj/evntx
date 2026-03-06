import { useRef } from "react"

type Props = {
  value: string
  onChange: (value: string) => void
  length?: number
}

function OTPInput({ value, onChange, length = 6 }: Props) {
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
    <div className="flex gap-3" onPaste={handlePaste}>
      {Array.from({ length }).map((_, index) => (
        <input
          key={index}
          ref={(el) => (inputs.current[index] = el)}
          maxLength={1}
          className="w-12 h-12 text-center border rounded-lg text-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          value={value[index] || ""}
          onChange={(e) => handleChange(index, e.target.value)}
          onKeyDown={(e) => handleKeyDown(index, e)}
        />
      ))}
    </div>
  )
}

export default OTPInput