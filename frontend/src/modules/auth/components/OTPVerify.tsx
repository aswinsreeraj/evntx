import { useState, useEffect } from "react"
import OTPInput from "./OTPInput"
import { authApi } from "../api"

export default function OTPVerify({ email }: any) {
  const [otp, setOtp] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  
  const [timer, setTimer] = useState(60)
  const [canResend, setCanResend] = useState(false)
  const [resending, setResending] = useState(false)

  useEffect(() => {
    let interval: ReturnType<typeof setInterval>
    if (timer > 0 && !canResend) {
      interval = setInterval(() => {
        setTimer((prev) => prev - 1)
      }, 1000)
    } else if (timer === 0) {
      setCanResend(true)
    }
    return () => clearInterval(interval)
  }, [timer, canResend])

  const verifyOtp = async () => {
    setError("")
    if (otp.length !== 6) {
      setError("Please enter a complete 6-digit OTP")
      return
    }
    
    setLoading(true)
    try {
      await authApi.verifyOtp(email, otp)
      window.location.href = "/profile"
    } catch (err: any) {
      console.error(err)
      setError(err.response?.data?.message || "Invalid or expired OTP. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  const handleResend = async () => {
    setError("")
    setResending(true)
    try {
      await authApi.requestOtp(email)
      setTimer(60)
      setCanResend(false)
      setOtp("")
    } catch (err: any) {
      console.error("Failed to resend OTP", err)
      setError("Failed to resend OTP. Please try again later.")
    } finally {
      setResending(false)
    }
  }

  return (
    <div className="flex flex-col items-center w-full max-w-sm m-auto">
      <h2 className="text-2xl font-bold mb-3 text-gray-900 text-center">
        Welcome back
      </h2>

      <p className="text-gray-500 mb-8 text-center text-sm leading-relaxed px-4">
        Dive back into the ultimate experience
      </p>

      <div className="w-full flex flex-col mb-6">
        <label className="text-sm font-medium text-gray-700 mb-2">Email</label>
        <div className="w-full border border-gray-200 bg-transparent rounded-xl px-4 py-3 text-sm text-gray-500 cursor-not-allowed">
          {email || "johnsmith@example.com"}
        </div>
      </div>

      <div className="w-full flex flex-col mb-2">
        <label className="text-sm font-medium text-gray-700 mb-2">
          One-Time Password
        </label>
        <OTPInput 
            value={otp} 
            onChange={(val: string) => {
                setOtp(val)
                setError("")
            }} 
        />
        {error && <p className="text-red-500 text-xs mt-2">{error}</p>}
      </div>

      <div className="w-full mb-8 flex items-center justify-between">
        <p className={`${canResend ? "text-gray-500" : "text-[#e53e5d]"} text-xs`}>
          OTP has been sent to the email
        </p>
        <button
            onClick={handleResend}
            disabled={!canResend || resending}
            className={`text-xs font-medium ${!canResend || resending ? "text-gray-400 cursor-not-allowed" : "text-blue-600 hover:text-blue-800 transition-colors"}`}
        >
            {resending ? "Resending..." : canResend ? "Resend OTP" : `Resend in ${timer}s`}
        </button>
      </div>

      <button
        onClick={verifyOtp}
        disabled={otp.length !== 6 || loading}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
      >
        {loading ? (
            <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
        ) : (
            "Login"
        )}
      </button>
    </div>
  )
}