import EventCarousel from "./EventCarousel";
import { useEvents } from "../../events/hooks";

export default function UpcomingEvents() {
  const { data } = useEvents({ sort: "date_asc", limit: 3 });
  const events = data?.events || [];

  return <EventCarousel title="Upcoming events" events={events} />;
}
