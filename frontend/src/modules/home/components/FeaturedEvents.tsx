import EventCarousel from "./EventCarousel";
import { useEvents } from "../../events/hooks";

export default function FeaturedEvents() {
  const { data } = useEvents({ sort: "created_at_desc", limit: 10 });
  const events = data?.events || [];

  return <EventCarousel title="Recently added events" events={events} />;
}
