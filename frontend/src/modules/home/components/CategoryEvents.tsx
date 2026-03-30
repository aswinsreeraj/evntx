import EventCarousel from "./EventCarousel";
import { useEvents } from "../../events/hooks";

interface CategoryEventsProps {
  activeCategory: string;
}

export default function CategoryEvents({ activeCategory }: CategoryEventsProps) {
  const { data } = useEvents({
    category: activeCategory !== "All" ? activeCategory : undefined,
    limit: 6,
  });
  const events = data?.events || [];

  return <EventCarousel title="" events={events} />;
}
