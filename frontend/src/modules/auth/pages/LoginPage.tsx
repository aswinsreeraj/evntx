import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { authApi } from "../api";
import { extractErrorMessage } from "../../../shared/utils/errorHandler";

function LoginPage() {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleRequestOtp = async () => {
    try {
      setLoading(true);
      setError("");

      await authApi.requestOtp(email);

      navigate("/verify-otp", { state: { email } });
    } catch (err: any) {
      setError(extractErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Login with Email OTP</h2>

      <input
        type="email"
        placeholder="Enter your email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />

      <button onClick={handleRequestOtp} disabled={loading}>
        {loading ? "Sending..." : "Send OTP"}
      </button>

      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
}

export default LoginPage;