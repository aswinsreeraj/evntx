import { useState } from "react";
import { X, ArrowLeft } from "lucide-react";
import Modal from "../../../shared/ui/Modal";
import OTPInput from "../../auth/components/OTPInput";
import { authApi } from "../../auth/api";
import { useNavigate } from "react-router-dom";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  email: string;
};

export default function AdminOtpModal({ isOpen, onClose, email }: Props) {
  const [otp, setOtp] = useState("");

  const navigate = useNavigate();

  const verifyOtp = async () => {
    if (otp.length !== 6) return;
    try {
      await authApi.verifyOtp(email, otp);
      onClose();
      navigate("/admin/users");
    } catch (e) {
      console.error(e);
      alert("Failed to verify OTP.");
    }
  };

  return (
    <Modal open={isOpen} onClose={onClose} className="p-0 max-w-4xl w-[90%] md:w-[800px] rounded-3xl overflow-hidden bg-white relative">
      <div className="flex flex-col md:flex-row h-full">
        {/* Left Side Image */}
        <div className="hidden md:block w-1/2 relative bg-gray-100">
          <img
            src="https://images.unsplash.com/photo-1552664730-d307ca884978?q=80&w=1470"
            alt="Admin Team"
            className="w-full h-full object-cover"
          />
        </div>

        {/* Right Side Content */}
        <div className="w-full md:w-1/2 p-10 flex flex-col justify-center relative min-h-[500px]">
          
          <button 
            onClick={onClose}
            className="absolute top-6 left-6 text-gray-400 hover:text-gray-600 transition"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>

          <button 
            onClick={onClose}
            className="absolute top-6 right-6 text-gray-400 hover:text-gray-600 transition"
          >
            <X className="w-5 h-5" />
          </button>

          <div className="flex flex-col w-full max-w-sm mx-auto mt-4">
            <h2 className="text-2xl font-bold mb-3 text-gray-900 text-center">
              Welcome back
            </h2>

            <p className="text-gray-500 mb-8 text-center text-sm leading-relaxed px-4">
              Let's go manage the platform.
            </p>

            <div className="w-full flex flex-col mb-6">
              <label className="text-sm font-medium text-gray-700 mb-2">Email</label>
              <div className="w-full border border-gray-200 bg-transparent rounded-xl px-4 py-3 text-sm text-gray-500 cursor-not-allowed">
                {email || "johnsmith@example.com"}
              </div>
            </div>

            <div className="w-full flex flex-col mb-2">
              <label className="text-sm font-medium text-gray-700 mb-2">
                One-Time Password
              </label>
              <OTPInput value={otp} onChange={setOtp} />
            </div>

            <div className="w-full mb-8">
              <p className="text-[#e53e5d] text-xs">
                OTP has been sent to the email
              </p>
            </div>

            <button
              onClick={verifyOtp}
              disabled={otp.length !== 6}
              className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Login
            </button>
          </div>
        </div>
      </div>
    </Modal>
  );
}
