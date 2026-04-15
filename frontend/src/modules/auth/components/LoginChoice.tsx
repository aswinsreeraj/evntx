import { useEffect, useState } from "react"
import { GoogleLogin } from "@react-oauth/google"
import { authApi } from "../api"

export default function LoginChoice({ setView, setEmail, isOrganizer, onClose }: any) {
  const [localEmail, setLocalEmail] = useState("")
  const [loading, setLoading] = useState(false)
  const [errorMsg, setErrorMsg] = useState("")
  const [allowGoogleLogin, setAllowGoogleLogin] = useState(true)

  useEffect(() => {
    const loadSettings = async () => {
      try {
        const settings = await authApi.getPlatformSettings();
        setAllowGoogleLogin(Boolean(settings?.allow_google_login));
      } catch {
        setAllowGoogleLogin(true);
      }
    };
    loadSettings();
  }, []);

  const handleGoogleSuccess = async (credentialResponse: any) => {
    setLoading(true);
    setErrorMsg("");
    try {
      if (credentialResponse.credential) {
        await authApi.googleLogin(credentialResponse.credential);
        if (typeof onClose === 'function') onClose();
        window.location.href = isOrganizer ? "/organizer/profile" : "/";
      } else {
        setErrorMsg("Google login failed. No credential received.");
      }
    } catch (error: any) {
      console.error("Google login failed", error);
      setErrorMsg(error.response?.data?.error?.message || "Google login failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleError = () => {
    console.error("Google Login Failed");
    setErrorMsg("Google Login Failed");
  };

  const handleContinue = async () => {
    if (!localEmail) return;
    setLoading(true);
    setErrorMsg("");
    try {
      const res = await authApi.requestOtp(localEmail);
      setEmail(localEmail);

      if (res.data?.is_new_user) {
        setView("register");
      } else {
        setView("otp-verify");
      }
    } catch (error: any) {
      console.error("Failed to send OTP", error);
      setErrorMsg(error.response?.data?.error?.message || "Failed to send OTP. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center w-full max-w-sm m-auto">
      <h2 className="text-2xl font-bold mb-3 text-gray-900 text-center">
        {isOrganizer ? "Manage Your Events" : "Welcome to the world of events"}
      </h2>

      <p className="text-gray-500 mb-8 text-center text-sm leading-relaxed px-4">
        {isOrganizer ? (
          <>Reach your audience and grow your community.<br/>Log in to EVNTX Organizer dashboard.</>
        ) : (
          <>Do you need to connect, learn, network or just chill?<br/>You can find an event here to get you there.</>
        )}
      </p>

      {allowGoogleLogin ? (
        <div className="w-full flex justify-center mb-8">
          <GoogleLogin
            onSuccess={handleGoogleSuccess}
            onError={handleGoogleError}
            theme="outline"
            size="large"
            text="continue_with"
          />
        </div>
      ) : (
        <p className="mb-8 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700">
          Google login is currently disabled by admin.
        </p>
      )}

      <div className="w-full flex items-center gap-4 mb-8">
        <div className="h-px bg-gray-200 flex-1"></div>
        <span className="text-xs text-gray-400 font-medium tracking-wide">or log in using your email</span>
        <div className="h-px bg-gray-200 flex-1"></div>
      </div>

      <input
        value={localEmail}
        placeholder="Enter your email here"
        className={`w-full border ${errorMsg ? 'border-red-500' : 'border-gray-300'} rounded-xl px-4 py-3 mb-2 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all text-sm`}
        onChange={(e) => { setLocalEmail(e.target.value); setErrorMsg(""); }}
      />
      {errorMsg ? (
        <p className="text-red-500 text-xs text-center w-full mb-6 font-medium bg-red-50 p-2 rounded-lg border border-red-100">{errorMsg}</p>
      ) : (
        <div className="mb-4"></div>
      )}

      <button
        onClick={handleContinue}
        disabled={!localEmail || loading}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors mb-8 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
      >
        {loading ? (
            <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
        ) : (
            "Continue"
        )}
      </button>

      <p className="text-xs text-center text-gray-500 leading-tight">
        By logging in, I agree to the <a href="#" className="text-indigo-900 hover:underline">Privacy Policy</a> and <br/>
        <a href="#" className="text-indigo-900 hover:underline">Terms of Service</a> of EVNTX.
      </p>
    </div>
  )
}
