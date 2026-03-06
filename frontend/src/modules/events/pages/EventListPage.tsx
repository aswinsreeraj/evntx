import { useEvents } from "../hooks";
import EventCard from "../components/EventCard";

function EventListPage() {
  const { data, isLoading, error } = useEvents();

  if (isLoading) return <p>Loading events...</p>;
  if (error) return <p>Failed to load events</p>;

  return (
    <div>
      <h2>Events</h2>

      {data.events.map((event: any) => (
        <EventCard key={event.id} event={event} />
      ))}
    </div>
  );
}

export default EventListPage;