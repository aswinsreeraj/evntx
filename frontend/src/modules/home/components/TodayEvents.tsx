import EventCarousel from "./EventCarousel";
import { useEvents } from "../../events/hooks";

export default function TodayEvents() {
  const { data } = useEvents({ limit: 9 })
  const events = data?.events?.slice(3, 6) || []

  return <EventCarousel title="Today's Events" events={events} />;
}
