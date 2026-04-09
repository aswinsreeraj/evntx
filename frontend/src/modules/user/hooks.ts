import { useQuery } from "@tanstack/react-query";
import { userApi } from "./api";
import { useAuthStore } from "../auth/store/authStore";

export const walletQueryKey = ["user-wallet"] as const;
export const walletTransactionsQueryKey = ["user-wallet-transactions"] as const;

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

export function useWalletTransactions(
  params: {
    page: number;
    limit: number;
    type?: "cr" | "dr";
    status?: "pending" | "completed" | "failed";
  },
  options?: { enabled?: boolean }
) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const userID = useAuthStore((state) => state.user?.id ?? null);
  const isEnabled = options?.enabled ?? true;

  return useQuery({
    queryKey: [
      ...walletTransactionsQueryKey,
      userID,
      params.page,
      params.limit,
      params.type ?? "",
      params.status ?? "",
    ],
    queryFn: () => userApi.getWalletTransactions(params),
    enabled: isAuthenticated && Boolean(userID) && isEnabled,
    staleTime: 30 * 1000,
  });
}

export function usePaymentSettings() {
  return useQuery({
    queryKey: ["payment-settings"],
    queryFn: () => userApi.getPaymentSettings(),
  });
}
