import { Navigate } from "react-router-dom";
import { useAuthStore } from "../../modules/auth/store/authStore";

type Props = {
  children: React.ReactNode;
  roles?: string[];
};

function ProtectedRoute({ children, roles }: Props) {
  const { isAuthenticated, roles: userRoles } = useAuthStore();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (roles && !roles.some((role) => userRoles.includes(role))) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

export default ProtectedRoute;