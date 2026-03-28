import { useQuery } from "@tanstack/react-query";
import { notificationsApi } from "./api";

export const notificationsQueryKey = (page: number, limit: number) =>
  ["notifications", page, limit] as const;

export function useNotifications(page = 1, limit = 6, enabled = true) {
  return useQuery({
    queryKey: notificationsQueryKey(page, limit),
    queryFn: () => notificationsApi.getNotifications(page, limit),
    enabled,
    refetchInterval: 45000,
  });
}
