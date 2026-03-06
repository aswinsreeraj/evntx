import { useState } from "react";
import { authApi } from "../../auth/api";

function AdminLoginPage() {
  const [email, setEmail] = useState("");

  const handleLogin = async () => {
    await authApi.requestOtp(email);
  };

  return (
    <div>
      <h2>Admin Login</h2>

      <input
        type="email"
        placeholder="Admin Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />

      <button onClick={handleLogin}>Send OTP</button>
    </div>
  );
}

export default AdminLoginPage;