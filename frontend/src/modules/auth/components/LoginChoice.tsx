import { handleGoogleLogin } from "../googleAuth"

export default function LoginChoice({ setView, setEmail }: any) {

  const handleGoogle = async () => {
    // later integrate Google SDK
    console.log("Google login")
  }

  return (
    <div className="flex flex-col items-center w-full max-w-sm mx-auto">
      <h2 className="text-2xl font-bold mb-3 text-gray-900 text-center">
        Welcome to the world of events
      </h2>

      <p className="text-gray-500 mb-8 text-center text-sm leading-relaxed px-4">
        Do you need to connect, learn, network or just chill?<br/>
        You can find an event here to get you there.
      </p>

      <button
        onClick={handleGoogle}
        className="w-full flex items-center justify-center gap-3 border border-gray-300 rounded-xl py-3 mb-8 hover:bg-gray-50 transition-colors"
      >
        <img 
          src="https://www.svgrepo.com/show/475656/google-color.svg" 
          alt="Google logo" 
          className="w-5 h-5"
        />
        <span className="text-gray-700 font-semibold text-sm">Login with Google</span>
      </button>

      <div className="w-full flex items-center gap-4 mb-8">
        <div className="h-px bg-gray-200 flex-1"></div>
        <span className="text-xs text-gray-400 font-medium tracking-wide">or log in using your email</span>
        <div className="h-px bg-gray-200 flex-1"></div>
      </div>

      <input
        placeholder="Enter your email here"
        className="w-full border border-gray-300 rounded-xl px-4 py-3 mb-6 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all text-sm"
        onChange={(e) => setEmail(e.target.value)}
      />

      <button
        onClick={() => setView("email-input")}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors mb-8"
      >
        Continue
      </button>

      <p className="text-xs text-center text-gray-500 leading-tight">
        By logging in, I agree to the <a href="#" className="text-indigo-900 hover:underline">Privacy Policy</a> and <br/>
        <a href="#" className="text-indigo-900 hover:underline">Terms of Service</a> of EVNTX.
      </p>
    </div>
  )
}