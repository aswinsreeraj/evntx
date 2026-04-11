import { useEffect } from "react";
import { useLocation } from "react-router-dom";
import { useEngagement } from "../shared/hooks/useEngagement";

export default function EngagementHydrator() {
  const location = useLocation();
  const { initialize, trackEvent } = useEngagement();

  useEffect(() => {
    initialize();
  }, [initialize]);

  useEffect(() => {
    trackEvent('page_view', undefined, JSON.stringify({ path: location.pathname }));
  }, [location.pathname, trackEvent]);

  return null;
}
