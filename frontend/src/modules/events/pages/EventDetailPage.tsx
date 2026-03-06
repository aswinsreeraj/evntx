import { useParams } from "react-router-dom";
import { useEvent } from "../hooks";

function EventDetailPage() {
  const { eventId } = useParams();

  const { data, isLoading, error } = useEvent(eventId!);

  if (isLoading) return <p>Loading event...</p>;
  if (error) return <p>Failed to load event</p>;

  return (
    <div>
      <h2>{data.title}</h2>

      <img
        src={data.cover_image_url}
        alt={data.title}
        style={{ width: "100%", maxHeight: "400px", objectFit: "cover" }}
      />

      <p>{data.description}</p>

      <h3>Tickets</h3>

      {data.ticket_types.map((ticket: any) => (
        <div key={ticket.id}>
          <p>
            {ticket.name} — ₹{ticket.price}
          </p>
        </div>
      ))}
    </div>
  );
}

export default EventDetailPage;