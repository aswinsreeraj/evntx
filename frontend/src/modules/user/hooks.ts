import { useQuery } from "@tanstack/react-query";
import { userApi } from "./api";
import { useAuthStore } from "../auth/store/authStore";

export const walletQueryKey = ["user-wallet"] as const;

export function useWallet(options?: { enabled?: boolean }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const userID = useAuthStore((state) => state.user?.id ?? null);
  const isEnabled = options?.enabled ?? true;

  return useQuery({
    queryKey: [...walletQueryKey, userID],
    queryFn: () => userApi.getWallet(),
    enabled: isAuthenticated && Boolean(userID) && isEnabled,
    staleTime: 60 * 1000,
  });
}
