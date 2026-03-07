import { useState } from "react"
import OTPInput from "./OTPInput"
import { authApi } from "../api"

export default function OTPVerify({ email }: any) {
  const [otp, setOtp] = useState("")

  const verifyOtp = async () => {
    if (otp.length !== 6) return
    try {
      await authApi.verifyOtp(email, otp)
      window.location.href = "/profile"
    } catch (e) {
      console.error(e)
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
        <OTPInput value={otp} onChange={setOtp} />
      </div>

      <div className="w-full mb-8">
        <p className="text-red-500 text-xs">
          OTP has been sent to the email
        </p>
      </div>

      <button
        onClick={verifyOtp}
        disabled={otp.length !== 6}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Login
      </button>
    </div>
  )
}