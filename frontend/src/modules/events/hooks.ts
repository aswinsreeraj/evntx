import { useQuery } from "@tanstack/react-query";
import { eventsApi } from "./api";

export function useEvents(params?: any) {
  return useQuery({
    queryKey: ["events", params],
    queryFn: () => eventsApi.listEvents(params),
  });
}

export function useEvent(eventId: string, isOrganizer = false, isAdmin = false) {
  return useQuery({
    queryKey: ["event", eventId, isOrganizer, isAdmin],
    queryFn: () => eventsApi.getEventForViewer(eventId, isOrganizer, isAdmin),
    enabled: Boolean(eventId),
    retry: false,
  });
}
