import { handleGoogleLogin } from "../googleAuth"

function LoginChoice({ setView, setEmail }: any) {

  const handleGoogle = async () => {
    // later integrate Google SDK
    console.log("Google login")
  }

  return (
    <div>
      <h2 className="text-xl font-semibold mb-2">
        Welcome to the world of events
      </h2>

      <p className="text-gray-500 mb-6">
        Connect, learn, network or just chill.
      </p>

      <button
        onClick={handleGoogle}
        className="w-full border rounded-lg py-3 mb-4"
      >
        Login with Google
      </button>

      <div className="text-center text-sm text-gray-400 mb-4">
        or log in using your email
      </div>

      <input
        placeholder="Enter your email"
        className="w-full border rounded-lg p-3 mb-4"
        onChange={(e) => setEmail(e.target.value)}
      />

      <button
        onClick={() => setView("email-input")}
        className="w-full bg-black text-white py-3 rounded-lg"
      >
        Continue
      </button>
    </div>
  )
}

export default LoginChoice