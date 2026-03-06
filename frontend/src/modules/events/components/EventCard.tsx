import { Link } from "react-router-dom";

type Props = {
  event: any;
};

function EventCard({ event }: Props) {
  return (
    <Link to={`/events/${event.id}`}>
      <div style={{ border: "1px solid #ccc", padding: "1rem", marginBottom: "1rem" }}>
        <img
          src={event.cover_image_url}
          alt={event.title}
          style={{ width: "100%", height: "200px", objectFit: "cover" }}
        />

        <h3>{event.title}</h3>
        <p>{event.city}</p>
        <p>{new Date(event.start_time).toLocaleString()}</p>
      </div>
    </Link>
  );
}

export default EventCard;