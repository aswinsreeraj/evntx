import { useCallback } from 'react';
import type { EngagementEventType } from '../api/engagement';
import { engagementApi } from '../api/engagement';
import { useAuthStore } from '../../modules/auth/store/authStore';

const SESSION_KEY = 'evntx_session_id';

let initializationPromise: Promise<string | null> | null = null;

const ensureSession = async (): Promise<string | null> => {
  let sessionId = sessionStorage.getItem(SESSION_KEY);
  if (sessionId) return sessionId;

  if (initializationPromise) return initializationPromise;

  initializationPromise = (async () => {
    try {
      const session = await engagementApi.initializeSession();
      sessionStorage.setItem(SESSION_KEY, session.id);
      return session.id;
    } catch {
      return null;
    } finally {
      initializationPromise = null;
    }
  })();

  return initializationPromise;
};

export const useEngagement = () => {
  const { roles, isAuthenticated } = useAuthStore();

  const initialize = useCallback(async () => {
    await ensureSession();
  }, []);

  const trackEvent = useCallback(async (type: EngagementEventType, eventId?: string, metadata?: string) => {
    
    if (isAuthenticated && (roles.includes('admin') || roles.includes('organizer'))) {
      return;
    }

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
