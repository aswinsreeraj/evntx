import { Outlet } from "react-router-dom";
import { useAuthStore } from "../modules/auth/store/authStore";
import { authApi } from "../modules/auth/api";

function Layout() {
  const { isAuthenticated } = useAuthStore();

  const handleLogout = async () => {
    await authApi.logout();
  };

  return (
    <div>
      <header style={{ padding: "1rem", borderBottom: "1px solid #ccc" }}>
        EVNTX
        {isAuthenticated && (
          <button onClick={handleLogout} style={{ marginLeft: "1rem" }}>
            Logout
          </button>
        )}
      </header>

      <main style={{ padding: "1rem" }}>
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;