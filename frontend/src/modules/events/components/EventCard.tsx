import { Link } from "react-router-dom"

type Props = {
  event: any
}

function EventCard({ event }: Props) {
  return (
    <Link to={`/events/${event.id}`}>
      <div className="bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-md transition">

        <div className="relative">
          <img
            src={event.cover_image_url}
            alt={event.title}
            className="h-48 w-full object-cover"
          />

          <span className="absolute top-3 right-3 bg-red-500 text-white text-xs px-3 py-1 rounded-full">
            {new Date(event.start_time).toLocaleDateString()}
          </span>
        </div>

        <div className="p-4">
          <h3 className="font-semibold">{event.title}</h3>

          <p className="text-sm text-gray-500">
            {new Date(event.start_time).toLocaleTimeString()}
          </p>

          <p className="text-sm text-red-500 mt-1">
            {event.city}
          </p>

          <p className="text-xs text-gray-400 mt-2">
            {event.description?.slice(0, 80)}
          </p>
        </div>
      </div>
    </Link>
  )
}

export default EventCard