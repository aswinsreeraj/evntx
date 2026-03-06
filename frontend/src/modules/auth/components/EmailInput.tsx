import { authApi } from "../api"

function EmailInput({ email, setEmail, setView }: any) {

  const requestOtp = async () => {
    try {
      await authApi.requestOtp(email)
      setView("otp-verify")
    } catch (error) {
      console.error("Failed to send OTP", error)
      alert("Failed to send OTP. Is your backend running and configured in .env?")
    }
  }

  return (
    <div>

      <h2 className="text-xl font-semibold mb-4">
        Welcome back
      </h2>

      <label>Email</label>

      <input
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        className="w-full border rounded-lg p-3 mb-6"
      />

      <button
        onClick={requestOtp}
        className="w-full bg-black text-white py-3 rounded-lg"
      >
        Send OTP +
      </button>

    </div>
  )
}

export default EmailInput