import { useEffect, useRef } from "react";
import { useAuthStore } from "../modules/auth/store/authStore";
import { userApi } from "../modules/user/api";

export default function AuthHydrator() {
  const { isAuthenticated, roles, user, setAuth, logout } = useAuthStore();
  const hydratedUserIdRef = useRef<string | null>(null);

  useEffect(() => {
    if (!isAuthenticated) {
      hydratedUserIdRef.current = null;
      return;
    }

    const currentUserId = user?.id ?? "__anonymous";
    if (hydratedUserIdRef.current === currentUserId) {
      return;
    }

    let cancelled = false;
    hydratedUserIdRef.current = currentUserId;

    const hydrateRoles = async () => {
      try {
        const profile = await userApi.getProfile();
        if (cancelled) return;

        setAuth(
          {
            id: profile.id,
            name: profile.name,
          },
          profile.roles || roles,
        );
      } catch {
        if (!cancelled && !user) {
          logout();
        }
      }
    };

    void hydrateRoles();

    return () => {
      cancelled = true;
    };
  }, [isAuthenticated, logout, roles, setAuth, user]);

  return null;
}
