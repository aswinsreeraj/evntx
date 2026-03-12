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
  
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);

  const validate = () => {
    const newErrors: Record<string, string> = {};
    const nameRegex = /^[a-zA-Z\s]+$/;

    if (!firstName.trim() || firstName.trim().length < 2) {
      newErrors.firstName = "First name must be at least 2 characters";
    } else if (!nameRegex.test(firstName)) {
      newErrors.firstName = "First name can only contain alphabets and spaces";
    }

    if (lastName.trim()) {
      if (lastName.trim().length < 2) {
        newErrors.lastName = "Last name must be at least 2 characters";
      } else if (!nameRegex.test(lastName)) {
        newErrors.lastName = "Last name can only contain alphabets and spaces";
      }
    }
    if (!dob) {
      newErrors.dob = "Date of birth is required";
    } else if (new Date(dob) > new Date()) {
      newErrors.dob = "Date of birth cannot be in the future";
    }
    if (!gender) {
      newErrors.gender = "Gender is required";
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleRegister = async () => {
    if (otp.length !== 6) {
      setErrors({ ...errors, api: "Please enter a valid 6-digit OTP" });
      return;
    }
    
    if (!validate()) return;
    
    setSubmitting(true);
    setErrors({});
    try {
      const name = `${firstName} ${lastName}`.trim();
      await authApi.register(email, otp, name, dob, gender);
      if (onClose) onClose();
      window.location.href = "/profile";
    } catch (e: any) {
      console.error(e);
      setErrors({ api: e.response?.data?.message || "Failed to register. Invalid OTP or request." });
    } finally {
      setSubmitting(false);
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


      <div className="w-full flex flex-col mb-4">
        <label className="text-sm font-medium text-gray-700 mb-2">Email</label>
        <div className="w-full border border-gray-200 bg-transparent rounded-xl px-4 py-3 text-sm text-gray-500 cursor-not-allowed">
          {email || "johnsmith@example.com"}
        </div>
      </div>


      <div className="w-full grid grid-cols-2 gap-4 mb-4">
        <div className="flex flex-col">
          <label className="text-sm font-medium text-gray-700 mb-2">First Name</label>
          <input
            value={firstName}
            onChange={(e) => { setFirstName(e.target.value); setErrors({ ...errors, firstName: "" }) }}
            placeholder="John"
            className={`w-full border ${errors.firstName ? 'border-red-500' : 'border-gray-300'} rounded-xl px-4 py-3 focus:outline-none focus:ring-1 ${errors.firstName ? 'focus:ring-red-400' : 'focus:ring-gray-400'} transition-colors text-sm`}
          />
          {errors.firstName && <span className="text-red-500 text-xs mt-1">{errors.firstName}</span>}
        </div>
        <div className="flex flex-col">
          <label className="text-sm font-medium text-gray-700 mb-2">Last Name</label>
          <input
            value={lastName}
            onChange={(e) => { setLastName(e.target.value); setErrors({ ...errors, lastName: "" }) }}
            placeholder="Smith"
            className={`w-full border ${errors.lastName ? 'border-red-500' : 'border-gray-300'} rounded-xl px-4 py-3 focus:outline-none focus:ring-1 ${errors.lastName ? 'focus:ring-red-400' : 'focus:ring-gray-400'} transition-colors text-sm`}
          />
          {errors.lastName && <span className="text-red-500 text-xs mt-1">{errors.lastName}</span>}
        </div>
      </div>


      <div className="w-full flex flex-col mb-4">
        <label className="text-sm font-medium text-gray-700 mb-2">Date of Birth</label>
        <div className="relative">
          <input
            type="date"
            value={dob}
            onChange={(e) => { setDob(e.target.value); setErrors({ ...errors, dob: "" }) }}
            className={`w-full border ${errors.dob ? 'border-red-500' : 'border-gray-300'} rounded-xl px-4 py-3 focus:outline-none focus:ring-1 ${errors.dob ? 'focus:ring-red-400' : 'focus:ring-gray-400'} transition-colors text-sm`}
          />
        </div>
        {errors.dob && <span className="text-red-500 text-xs mt-1">{errors.dob}</span>}
      </div>


      <div className="w-full flex flex-col mb-4 relative">
        <label className="text-sm font-medium text-gray-700 mb-2">Gender</label>
        <div className="grid grid-cols-3 gap-3">
          <button
            onClick={() => { setGender("Male"); setErrors({ ...errors, gender: "" }) }}
            className={`py-3 rounded-xl border text-sm font-medium transition-colors ${
              gender === "Male" 
                ? "bg-[#e53e5d] text-white border-transparent" 
                : "bg-white text-gray-600 border-gray-300 hover:bg-gray-50"
            }`}
          >
            Male
          </button>
          
          <button
            onClick={() => { setGender("Female"); setErrors({ ...errors, gender: "" }) }}
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

            {isOtherDropdownOpen && (
              <div className="absolute top-full left-0 mt-1 w-full bg-white border border-gray-200 rounded-lg shadow-lg z-10 overflow-hidden">
                <button 
                  onClick={() => { setGender("Other"); setIsOtherDropdownOpen(false); setErrors({ ...errors, gender: "" }) }}
                  className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  Other
                </button>
              </div>
            )}
          </div>
        </div>
        {errors.gender && <span className="text-red-500 text-xs mt-1">{errors.gender}</span>}
      </div>


      <div className="w-full flex flex-col mb-8 mt-2">
        <label className="text-sm font-medium text-gray-700 mb-2">Verification Code</label>
        <OTPInput value={otp} onChange={(val: string) => { setOtp(val); setErrors({...errors, api: ""}) }} />
        <p className="text-[#e53e5d] text-xs mt-2">
          Verification code has been sent to the email
        </p>
      </div>

      {errors.api && (
        <div className="w-full mb-4 p-3 bg-red-50 border border-red-100 rounded-lg">
          <p className="text-red-600 text-sm text-center font-medium">{errors.api}</p>
        </div>
      )}

      <button 
        onClick={handleRegister}
        disabled={otp.length !== 6 || submitting}
        className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors mt-2 disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
      >
        {submitting ? (
          <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
        ) : (
          "Register"
        )}
      </button>

    </div>
  )
}