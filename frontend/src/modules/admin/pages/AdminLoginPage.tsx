import { useState } from "react";
import { authApi } from "../../auth/api";
import AdminOtpModal from "../components/AdminOtpModal";

export default function AdminLoginPage() {
  const [email, setEmail] = useState("");
  const [isOtpModalOpen, setIsOtpModalOpen] = useState(false);

  const handleLogin = async () => {
    if (!email) return;
    try {
      await authApi.requestOtp(email);
      setIsOtpModalOpen(true);
    } catch (error) {
      console.error("Failed to send OTP", error);
      alert("Failed to send OTP. Is backend running?");
    }
  };

  return (
    <div className="min-h-screen flex bg-white">

      <div className="hidden lg:flex w-[20%] bg-[#0b101e] flex-col justify-start p-10">
        <div className="flex flex-col items-start gap-1">
          <h1 className="text-white text-2xl font-black tracking-wider">EVNTX</h1>
          <p className="text-gray-400 text-xs">Admin Panel</p>
        </div>
      </div>


      <div className="flex-1 flex flex-col items-center justify-center p-8">
        <div className="w-full max-w-md flex flex-col items-center">
          
          <h2 className="text-3xl font-bold text-gray-900 mb-3 text-center">
            Welcome to Platform Management
          </h2>
          
          <p className="text-gray-500 mb-10 text-center text-sm">
            Manage platform settings as administrator
          </p>

          <button className="w-full bg-white border border-gray-200 text-gray-700 py-3.5 rounded-xl flex items-center justify-center gap-3 font-medium text-sm hover:bg-gray-50 mb-8 transition-colors">
            <svg viewBox="0 0 24 24" className="w-5 h-5" aria-hidden="true">
              <path d="M12.0003 4.75C13.7703 4.75 15.3553 5.36002 16.6053 6.54998L20.0303 3.125C17.9502 1.19 15.2353 0 12.0003 0C7.31028 0 3.25527 2.69 1.28027 6.60998L5.27028 9.70498C6.21525 6.81002 8.87028 4.75 12.0003 4.75Z" fill="#EA4335"/>
              <path d="M23.49 12.275C23.49 11.49 23.415 10.73 23.3 10H12V14.51H18.47C18.18 15.99 17.34 17.25 16.08 18.1L19.945 21.1C22.2 19.01 23.49 15.92 23.49 12.275Z" fill="#4285F4"/>
              <path d="M5.26498 14.2949C5.02498 13.5699 4.88501 12.7999 4.88501 11.9999C4.88501 11.1999 5.01998 10.4299 5.26498 9.7049L1.275 6.60986C0.46 8.22986 0 10.0599 0 11.9999C0 13.9399 0.46 15.7699 1.28 17.3899L5.26498 14.2949Z" fill="#FBBC05"/>
              <path d="M12.0004 24.0001C15.2404 24.0001 17.9654 22.935 19.9454 21.095L16.0804 18.095C15.0054 18.82 13.6204 19.245 12.0004 19.245C8.8704 19.245 6.21537 17.185 5.26538 14.29L1.27539 17.385C3.25539 21.31 7.3104 24.0001 12.0004 24.0001Z" fill="#34A853"/>
            </svg>
            Login with Google
          </button>

          <div className="w-full flex items-center gap-4 mb-8">
            <div className="flex-1 h-px bg-gray-200"></div>
            <span className="text-gray-400 text-sm">or log in using your email</span>
            <div className="flex-1 h-px bg-gray-200"></div>
          </div>

          <div className="w-full mb-6 relative">
            <input
              type="email"
              placeholder="Enter your email here"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full border border-gray-300 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all placeholder-gray-400"
            />
          </div>

          <button 
            onClick={handleLogin}
            disabled={!email}
            className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Continue
          </button>
        </div>
      </div>

      <AdminOtpModal 
        isOpen={isOtpModalOpen} 
        onClose={() => setIsOtpModalOpen(false)} 
        email={email}
      />
    </div>
  );
}