import { useLocation, useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import { authApi } from "../api";
import { extractErrorMessage } from "../../../shared/utils/errorHandler";

function VerifyOtpPage() {
  const location = useLocation();
  const navigate = useNavigate();

  const email = location.state?.email;

  const [otp, setOtp] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [timer, setTimer] = useState(60);

  useEffect(() => {
    if (!email) {
      navigate("/login");
    }
  }, [email, navigate]);

  useEffect(() => {
    if (timer <= 0) return;

    const interval = setInterval(() => {
      setTimer((prev) => prev - 1);
    }, 1000);

    return () => clearInterval(interval);
  }, [timer]);

  const handleVerify = async () => {
    try {
      setLoading(true);
      setError("");

      await authApi.verifyOtp(email, otp);

      navigate("/");
    } catch (err: any) {
      setError(extractErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleResend = async () => {
    if (timer > 0) return;

    try {
      await authApi.requestOtp(email);
      setTimer(60);
    } catch (err: any) {
      setError(extractErrorMessage(err));
    }
  };

  return (
    <div>
      <h2>Verify OTP</h2>
      <p>Email: {email}</p>

      <input
        type="text"
        placeholder="Enter OTP"
        value={otp}
        onChange={(e) => setOtp(e.target.value)}
      />

      <button onClick={handleVerify} disabled={loading}>
        {loading ? "Verifying..." : "Verify"}
      </button>

      <div>
        {timer > 0 ? (
          <p>Resend OTP in {timer}s</p>
        ) : (
          <button onClick={handleResend}>Resend OTP</button>
        )}
      </div>

      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
}

export default VerifyOtpPage;