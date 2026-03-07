import { authApi } from "../api"

export default function EmailInput({ email, setEmail, setView }: any) {

  const requestOtp = async () => {
    try {
      if (!email) return;
      const res = await authApi.requestOtp(email)
      
      if (res.data?.is_new_user) {
        setView("register")
      } else {
        setView("otp-verify")
      }
    } catch (error) {
      console.error("Failed to send OTP", error)
      alert("Failed to send OTP. Is your backend running and configured in .env?")
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
          onChange={(e) => setEmail(e.target.value)}
          placeholder="johnsmith@example.com"
          className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all text-sm"
        />
      </div>

      <button
        onClick={requestOtp}
        disabled={!email}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Continue
      </button>
    </div>
  )
}