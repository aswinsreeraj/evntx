import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { adminApi } from "./api";

export function useUsers(params: any) {
  return useQuery({
    queryKey: ["admin-users", params],
    queryFn: () => adminApi.getUsers(params),
  });
}

export function useToggleUserStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, isActive }: any) =>
      adminApi.updateUserStatus(userId, isActive),

    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
      queryClient.invalidateQueries({ queryKey: ["admin-organizers"] });
    },
  });
}

export function useOrganizers(params: any) {
  return useQuery({
    queryKey: ["admin-organizers", params],
    queryFn: () => adminApi.getOrganizers(params),
  });
}

export function useApproveOrganizer() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (organizerId: string) => adminApi.approveOrganizer(organizerId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-organizers"] });
    },
  });
}

export function useRejectOrganizer() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (organizerId: string) => adminApi.rejectOrganizer(organizerId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-organizers"] });
    },
  });
}

export function useEvents(params: any) {
  return useQuery({
    queryKey: ["admin-events", params],
    queryFn: () => adminApi.getEvents(params),
  });
}

export function useApproveEvent() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (eventId: string) => adminApi.approveEvent(eventId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-events"] });
    },
  });
}

export function useRejectEvent() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ eventId, reason }: { eventId: string; reason: string }) =>
      adminApi.rejectEvent(eventId, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-events"] });
    },
  });
}

export function usePlatformWallet() {
  return useQuery({
    queryKey: ["admin-platform-wallet"],
    queryFn: () => adminApi.getPlatformWallet(),
  });
}

export function useSettings() {
  return useQuery({
    queryKey: ["admin-settings"],
    queryFn: () => adminApi.getPlatformSettings(),
  });
}

export function useUpdateSettings() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Partial<Parameters<typeof adminApi.updatePlatformSettings>[0]>) =>
      adminApi.updatePlatformSettings(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-settings"] });
    },
  });
}

export function usePaymentSettings() {
  return useQuery({
    queryKey: ["admin-payment-settings"],
    queryFn: () => adminApi.getPaymentSettings(),
  });
}

export function useUpdatePaymentProvider() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ provider, data }: { provider: string; data: { is_enabled: boolean; config: Record<string, any> } }) =>
      adminApi.updatePaymentProvider(provider, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-payment-settings"] });
    },
  });
}

export function useAdmins() {
  return useQuery({
    queryKey: ["admin-users-list"], 
    queryFn: () => adminApi.getAdmins(),
  });
}

export function useAuditLogs(page: number, limit: number) {
  return useQuery({
    queryKey: ["admin-audit-logs", page, limit],
    queryFn: () => adminApi.getAuditLogs(page, limit),
  });
}
