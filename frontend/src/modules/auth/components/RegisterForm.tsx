import { useState } from "react";
import { ChevronDown, CalendarDays } from "lucide-react";

export default function RegisterForm({ email }: any) {
  const [gender, setGender] = useState<"Male" | "Female" | "Other" | "">("Male");
  const [isOtherDropdownOpen, setIsOtherDropdownOpen] = useState(false);

  return (
    <div className="flex flex-col items-center w-full max-w-sm mx-auto overflow-y-auto max-h-[80vh] pr-2 pb-4 scrollbar-hide" style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}>
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
            placeholder="John"
            className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-1 focus:ring-gray-400 transition-colors text-sm"
          />
        </div>
        <div className="flex flex-col">
          <label className="text-sm font-medium text-gray-700 mb-2">Last Name</label>
          <input
            placeholder="John"
            className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-1 focus:ring-gray-400 transition-colors text-sm"
          />
        </div>
      </div>

      {/* Date of Birth */}
      <div className="w-full flex flex-col mb-4">
        <label className="text-sm font-medium text-gray-700 mb-2">Date of Birth</label>
        <div className="relative">
          <input
            type="text"
            placeholder="Select date"
            className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-1 focus:ring-gray-400 transition-colors text-sm"
          />
          <CalendarDays className="absolute right-4 top-1/2 -translate-y-1/2 text-red-500 w-5 h-5 pointer-events-none" />
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
      <div className="w-full flex flex-col mb-8">
        <label className="text-sm font-medium text-gray-700 mb-2">Verification Code</label>
        <input
          placeholder="Enter verification code"
          className="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-1 focus:ring-gray-400 transition-colors text-sm mb-2"
        />
        <p className="text-[#e53e5d] text-xs">
          Verification code has been sent to the email
        </p>
      </div>

      {/* Register Button */}
      <button className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors mt-2">
        Register
      </button>

    </div>
  )
}