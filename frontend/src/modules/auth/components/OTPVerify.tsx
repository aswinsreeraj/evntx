import { useState } from "react"
import OTPInput from "./OTPInput"
import { authApi } from "../api"

function OTPVerify({ email }: any) {
  const [otp, setOtp] = useState("")

  const verifyOtp = async () => {
    if (otp.length !== 6) return
    await authApi.verifyOtp(email, otp)
  }

  return (
    <div>

      <h2 className="text-xl font-semibold mb-4">
        Welcome back
      </h2>

      <label className="text-sm text-gray-500">
        Email
      </label>

      <input
        value={email}
        disabled
        className="w-full border rounded-lg p-3 mb-6"
      />

      <label className="text-sm text-gray-500 mb-2 block">
        One-Time Password
      </label>

      <OTPInput value={otp} onChange={setOtp} />

      <p className="text-red-500 text-sm mt-3 mb-6">
        OTP has been sent to the email
      </p>

      <button
        onClick={verifyOtp}
        className="w-full bg-black text-white py-3 rounded-lg"
      >
        Login
      </button>

    </div>
  )
}

export default OTPVerify