import { Navigate, useLocation } from "react-router-dom";
import { useAuthStore } from "../../modules/auth/store/authStore";

type Props = {
  children: React.ReactNode;
  roles?: string[];
};

function ProtectedRoute({ children, roles }: Props) {
  const location = useLocation();
  const { isAuthenticated, roles: userRoles } = useAuthStore();
  const isAdminRoute = location.pathname.startsWith("/admin");
  const isNormalUser = !userRoles.includes("admin") && !userRoles.includes("organizer");

  if (!isAuthenticated) {
    return <Navigate to={isAdminRoute ? "/admin/login" : "/login"} replace />;
  }

  const hasRequiredRole = !roles
    ? true
    : roles.some((role) => {
        if (role === "goer") {
          return isNormalUser;
        }
        return userRoles.includes(role);
      });

  if (!hasRequiredRole) {
    const readableRoles = roles?.map((role) => (role === "goer" ? "normal user" : role)) ?? [];
    return (
      <div className="flex min-h-[60vh] items-center justify-center px-6 py-12">
        <div className="w-full max-w-xl rounded-[24px] border border-[#ececec] bg-white p-8 text-center shadow-[0_12px_32px_rgba(15,23,42,0.08)]">
          <div className="text-sm font-semibold uppercase tracking-[0.24em] text-[#ff445d]">Access Denied</div>
          <h1 className="mt-3 text-2xl font-semibold text-[#111827]">You are not authorized to view this page.</h1>
          <p className="mt-3 text-sm leading-6 text-[#6b7280]">
            This route is available only for {readableRoles.join(" / ")} accounts. Please use the correct dashboard for your role.
          </p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}

export default ProtectedRoute;
