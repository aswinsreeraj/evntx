import { useState, useEffect } from "react";
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
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  
  const [timer, setTimer] = useState(60);
  const [canResend, setCanResend] = useState(false);
  const [resending, setResending] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    let interval: ReturnType<typeof setInterval>;
    if (isOpen && timer > 0 && !canResend) {
      interval = setInterval(() => {
        setTimer((prev) => prev - 1);
      }, 1000);
    } else if (timer === 0) {
      setCanResend(true);
    }
    return () => clearInterval(interval);
  }, [isOpen, timer, canResend]);

  useEffect(() => {
    if (isOpen) {
      setOtp("");
      setError("");
      setTimer(60);
      setCanResend(false);
    }
  }, [isOpen]);

  const verifyOtp = async () => {
    setError("");
    if (otp.length !== 6) {
      setError("Please enter a complete 6-digit OTP");
      return;
    }
    
    setLoading(true);
    try {
      await authApi.verifyOtp(email, otp);
      onClose();
      navigate("/admin/users");
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.message || "Invalid or expired OTP. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const handleResend = async () => {
    setError("");
    setResending(true);
    try {
      await authApi.requestOtp(email);
      setTimer(60);
      setCanResend(false);
      setOtp("");
    } catch (err: any) {
      console.error("Failed to resend OTP", err);
      setError("Failed to resend OTP. Please try again later.");
    } finally {
      setResending(false);
    }
  };

  return (
    <Modal open={isOpen} onClose={onClose} className="p-0 max-w-4xl w-[90%] md:w-[800px] rounded-3xl overflow-hidden bg-white relative">
      <div className="flex flex-col md:flex-row h-full">

        <div className="hidden md:block w-1/2 relative bg-gray-100">
          <img
            src="https://images.unsplash.com/photo-1552664730-d307ca884978?q=80&w=1470"
            alt="Admin Team"
            className="w-full h-full object-cover"
          />
        </div>


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
              <OTPInput 
                value={otp} 
                onChange={(val: string) => {
                  setOtp(val);
                  setError("");
                }} 
              />
              {error && <p className="text-red-500 text-xs mt-2">{error}</p>}
            </div>

            <div className="w-full mb-8 flex items-center justify-between">
              <p className={`${canResend ? "text-gray-500" : "text-[#e53e5d]"} text-xs`}>
                OTP has been sent to the email
              </p>
              <button
                onClick={handleResend}
                disabled={!canResend || resending}
                className={`text-xs font-medium ${!canResend || resending ? "text-gray-400 cursor-not-allowed" : "text-blue-600 hover:text-blue-800 transition-colors"}`}
              >
                {resending ? "Resending..." : canResend ? "Resend OTP" : `Resend in ${timer}s`}
              </button>
            </div>

            <button
              onClick={verifyOtp}
              disabled={otp.length !== 6 || loading}
              className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl font-medium text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
            >
              {loading ? (
                <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
              ) : (
                "Login"
              )}
            </button>
          </div>
        </div>
      </div>
    </Modal>
  );
}
