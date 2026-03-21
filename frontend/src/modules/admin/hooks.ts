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