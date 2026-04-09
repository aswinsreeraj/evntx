import { useCallback } from 'react';
import type { EngagementEventType } from '../api/engagement';
import { engagementApi } from '../api/engagement';

const SESSION_KEY = 'evntx_session_id';

const ensureSession = async (): Promise<string | null> => {
  let sessionId = sessionStorage.getItem(SESSION_KEY);
  if (!sessionId) {
    try {
      const session = await engagementApi.initializeSession();
      sessionStorage.setItem(SESSION_KEY, session.id);
      sessionId = session.id;
    } catch {
      return null;
    }
  }
  return sessionId;
};

export const useEngagement = () => {
  const initialize = useCallback(async () => {
    await ensureSession();
  }, []);

  const trackEvent = useCallback(async (type: EngagementEventType, eventId?: string, metadata?: string) => {
    try {
      const sessionId = await ensureSession();
      if (!sessionId) return;

      await engagementApi.trackEvent({
        session_id: sessionId,
        event_type: type,
        event_id: eventId,
        metadata: metadata
      });
    } catch (error) {
      console.error(`Failed to track ${type} event:`, error);
    }
  }, []);

  return { initialize, trackEvent };
};
