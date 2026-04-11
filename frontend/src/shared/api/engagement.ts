import api from "../../services/axios";

export type EngagementEventType = 'page_view' | 'event_view' | 'ticket_selected' | 'checkout_started';

export interface TrackEventPayload {
  session_id: string;
  event_type: EngagementEventType;
  event_id?: string;
  metadata?: string;
}

export interface VisitorSession {
  id: string;
  user_id?: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
}

export const engagementApi = {
  initializeSession: async (): Promise<VisitorSession> => {
    const res = await api.post('/engagement/session');
    return res.data;
  },

  trackEvent: async (payload: TrackEventPayload): Promise<void> => {
    await api.post('/engagement/track', payload);
  }
};
