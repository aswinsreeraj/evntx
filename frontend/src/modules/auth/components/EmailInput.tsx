import { useState } from "react"
import { authApi } from "../api"

export default function EmailInput({ email, setEmail, setView }: any) {
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const validateEmail = (email: string) => {
    return String(email)
      .toLowerCase()
      .match(
        /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      );
  };

  const requestOtp = async () => {
    setError("")
    if (!email) {
      setError("Email is required")
      return;
    }
    if (!validateEmail(email)) {
      setError("Please enter a valid email address")
      return;
    }

    setLoading(true)
    try {
      const res = await authApi.requestOtp(email)
      
      if (res.data?.is_new_user) {
        setView("register")
      } else {
        setView("otp-verify")
      }
    } catch (err: any) {
      console.error("Failed to send OTP", err)
      setError(err.response?.data?.message || "Failed to send OTP. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col items-center w-full max-w-sm mx-auto">
      <h2 className="text-2xl font-bold mb-3 text-gray-900 text-center">
        Welcome back
      </h2>

      <p className="text-gray-500 mb-8 text-center text-sm leading-relaxed px-4">
        Dive back into the ultimate experience
      </p>

      <div className="w-full flex flex-col mb-8">
        <label className="text-sm font-medium text-gray-700 mb-2">Email</label>
        <input
          value={email}
          onChange={(e) => {
            setEmail(e.target.value)
            setError("")
          }}
          placeholder="johnsmith@example.com"
          className={`w-full border ${error ? 'border-red-500' : 'border-gray-300'} rounded-xl px-4 py-3 focus:outline-none focus:ring-2 ${error ? 'focus:ring-red-500/20 focus:border-red-500' : 'focus:ring-blue-500/20 focus:border-blue-500'} transition-all text-sm`}
        />
        {error && <p className="text-red-500 text-xs mt-2">{error}</p>}
      </div>

      <button
        onClick={requestOtp}
        disabled={!email || loading}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
      >
        {loading ? (
            <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
        ) : (
            "Continue"
        )}
      </button>
    </div>
  )
}