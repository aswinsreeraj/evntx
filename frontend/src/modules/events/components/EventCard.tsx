import { Link } from "react-router-dom"
import { MapPin } from "lucide-react"

type Props = {
  event: any
}

export default function EventCard({ event }: Props) {
  // Format date correctly based on Figma (e.g., 21 Feb 2026)
  const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { day: '2-digit', month: 'short', year: 'numeric' };
    return new Date(dateString).toLocaleDateString('en-GB', options);
  };

  // Format time (e.g., 12:00 PM)
  const formatTime = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { hour: '2-digit', minute: '2-digit' };
    return new Date(dateString).toLocaleTimeString('en-US', options);
  };

  return (
    <Link to={`/events/${event.id}`}>
      <div className="bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-shadow h-full flex flex-col border border-gray-100">

        <div className="relative h-48 w-full shrink-0">
          <img
            src={event.cover_image_url}
            alt={event.title}
            className="h-full w-full object-cover"
          />
          <span className="absolute top-3 right-3 bg-[#e53e5d]/90 backdrop-blur-sm text-white text-xs font-medium px-3 py-1.5 rounded-md shadow-sm">
            {formatDate(event.start_time)}
          </span>
        </div>

        <div className="p-5 flex flex-col flex-grow">
          <h3 className="font-bold text-gray-900 text-lg mb-1 line-clamp-1">{event.title}</h3>

          <p className="text-sm font-medium text-gray-700 mb-1">
            {formatTime(event.start_time)}
          </p>

          <div className="flex items-center text-sm text-[#e53e5d] font-medium mb-3">
            <MapPin className="w-3.5 h-3.5 mr-1" />
            {event.city}
          </div>

          <p className="text-xs text-gray-500 mb-4 line-clamp-2 leading-relaxed flex-grow">
            {event.description}
          </p>

          {/* Tags */}
          <div className="flex flex-wrap gap-2 mt-auto">
            {event.tags?.map((tag: string, index: number) => (
              <span 
                key={index} 
                className="bg-gray-100 text-gray-600 text-xs font-medium px-3 py-1 rounded-full"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
      </div>
    </Link>
  )
}