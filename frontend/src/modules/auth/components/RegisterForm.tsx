import { useState } from "react";
import { ChevronDown } from "lucide-react";
import OTPInput from "./OTPInput";
import { authApi } from "../api";

export default function RegisterForm({ email, onClose }: any) {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [dob, setDob] = useState("");
  const [gender, setGender] = useState<"Male" | "Female" | "Other" | "">("Male");
  const [isOtherDropdownOpen, setIsOtherDropdownOpen] = useState(false);
  const [otp, setOtp] = useState("");

  const handleRegister = async () => {
    if (otp.length !== 6 || !firstName) return;
    try {
      const name = `${firstName} ${lastName}`.trim();
      await authApi.verifyOtp(email, otp, name);
      if (onClose) onClose();
      window.location.href = "/profile";
    } catch (e) {
      console.error(e);
      alert("Failed to register. Invalid OTP.");
    }
  };

  return (
    <div className="flex flex-col items-center w-full max-w-sm m-auto py-4">
      <h2 className="text-2xl font-bold mb-3 text-gray-900 text-center mt-4">
        Welcome to EVNTX family
      </h2>

      <p className="text-gray-500 mb-6 text-center text-sm leading-relaxed px-4">
        Let's complete the registration
      </p>

      {/* Email Input (Disabled state) */}
      <div className="w-full flex flex-col mb-4">
        <label className="text-sm font-medium text-gray-700 mb-2">Email</label>
        <div className="w-full border border-gray-200 bg-transparent rounded-xl px-4 py-3 text-sm text-gray-500 cursor-not-allowed">
          {email || "johnsmith@example.com"}
        </div>
      </div>

      {/* Name Inputs */}
      <div className="w-full grid grid-cols-2 gap-4 mb-4">
        <div className="flex flex-col">
          <label className="text-sm font-medium text-gray-700 mb-2">First Name</label>
          <input
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            placeholder="John"
            className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-1 focus:ring-gray-400 transition-colors text-sm"
          />
        </div>
        <div className="flex flex-col">
          <label className="text-sm font-medium text-gray-700 mb-2">Last Name</label>
          <input
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            placeholder="Smith"
            className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-1 focus:ring-gray-400 transition-colors text-sm"
          />
        </div>
      </div>

      {/* Date of Birth */}
      <div className="w-full flex flex-col mb-4">
        <label className="text-sm font-medium text-gray-700 mb-2">Date of Birth</label>
        <div className="relative">
          <input
            type="date"
            value={dob}
            onChange={(e) => setDob(e.target.value)}
            className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-1 focus:ring-gray-400 transition-colors text-sm"
          />
        </div>
      </div>

      {/* Gender Selection */}
      <div className="w-full flex flex-col mb-4 relative">
        <label className="text-sm font-medium text-gray-700 mb-2">Gender</label>
        <div className="grid grid-cols-3 gap-3">
          <button
            onClick={() => setGender("Male")}
            className={`py-3 rounded-xl border text-sm font-medium transition-colors ${
              gender === "Male" 
                ? "bg-[#e53e5d] text-white border-transparent" 
                : "bg-white text-gray-600 border-gray-300 hover:bg-gray-50"
            }`}
          >
            Male
          </button>
          
          <button
            onClick={() => setGender("Female")}
            className={`py-3 rounded-xl border text-sm font-medium transition-colors ${
              gender === "Female" 
                ? "bg-[#e53e5d] text-white border-transparent" 
                : "bg-white text-gray-600 border-gray-300 hover:bg-gray-50"
            }`}
          >
            Female
          </button>
          
          <div className="relative">
            <button
              onClick={() => setIsOtherDropdownOpen(!isOtherDropdownOpen)}
              className={`w-full py-3 px-3 flex items-center justify-between rounded-xl border text-sm font-medium transition-colors ${
                gender === "Other" 
                  ? "bg-[#e53e5d] text-white border-transparent" 
                  : "bg-white text-gray-600 border-gray-300 hover:bg-gray-50"
              }`}
            >
              Other
              <ChevronDown className="w-4 h-4 ml-1 opacity-70" />
            </button>
            {/* Minimal Dropdown implementation */}
            {isOtherDropdownOpen && (
              <div className="absolute top-full left-0 mt-1 w-full bg-white border border-gray-200 rounded-lg shadow-lg z-10 overflow-hidden">
                <button 
                  onClick={() => { setGender("Other"); setIsOtherDropdownOpen(false); }}
                  className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  Other
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Verification Code */}
      <div className="w-full flex flex-col mb-8 mt-2">
        <label className="text-sm font-medium text-gray-700 mb-2">Verification Code</label>
        <OTPInput value={otp} onChange={setOtp} />
        <p className="text-[#e53e5d] text-xs mt-2">
          Verification code has been sent to the email
        </p>
      </div>

      {/* Register Button */}
      <button 
        onClick={handleRegister}
        disabled={otp.length !== 6 || !firstName}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors mt-2 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Register
      </button>

    </div>
  )
}