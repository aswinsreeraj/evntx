import { Link } from "react-router-dom"
import { MapPin } from "lucide-react"

type Props = {
  event: any
}

export default function EventCard({ event }: Props) {
  const title = event.title ?? event.Title
  const city = event.city ?? event.City
  const category = event.category ?? event.Category
  const venueName = event.venue_name ?? event.VenueName
  const coverImageUrl = event.cover_image_url ?? event.CoverImageURL
  const startTime = event.start_time ?? event.StartTime
  const description = event.description ?? event.Description
  const eventLink = event.slug ?? event.Slug ?? event.id ?? event.ID
  const rawTags = event.tags ?? event.Tags
  const eventTags = Array.isArray(rawTags)
    ? rawTags
    : typeof rawTags === "string" && rawTags.length > 0
      ? rawTags.split(",").map((tag: string) => tag.trim()).filter(Boolean)
      : category
        ? [category]
        : []

  const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { day: '2-digit', month: 'short', year: 'numeric' };
    return new Date(dateString).toLocaleDateString('en-GB', options);
  };

  const formatTime = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { hour: '2-digit', minute: '2-digit' };
    return new Date(dateString).toLocaleTimeString('en-US', options);
  };

  return (
    <Link to={`/events/${eventLink}`}>
      <div className=" h-98 w-98 bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-shadow flex flex-col border border-gray-100">

        <div className="relative h-48 w-full shrink-0">
          {coverImageUrl ? (
            <img
              src={coverImageUrl.startsWith("/") ? `${import.meta.env.VITE_API_BASE_URL}${coverImageUrl}` : coverImageUrl}
              alt={title}
              className="h-full w-full object-cover"
            />
          ) : (
            <div className="h-full w-full bg-gradient-to-br from-gray-700 to-gray-900 flex items-center justify-center">
              <span className="text-white text-4xl font-bold opacity-30">{(title || "E")[0].toUpperCase()}</span>
            </div>
          )}
          <span className="absolute top-3 right-3 bg-[#E7364D]/60 backdrop-blur-xl text-white text-xs font-medium px-3 py-1.5 rounded-full">
            {startTime ? formatDate(startTime) : "Live"}
          </span>
        </div>

        <div className="p-5 flex flex-col flex-grow">
          <h3 className="font-bold text-gray-900 text-lg mb-1 line-clamp-1">{title}</h3>

          <p className="text-sm font-medium text-gray-700 mb-1">
            {startTime ? formatTime(startTime) : "See details"}
          </p>

          <div className="flex items-center text-sm text-[#e53e5d] font-medium mb-3">
            <MapPin className="w-3.5 h-3.5 mr-1" />
            {city}
          </div>

          <p className="text-xs text-gray-500 mb-4 line-clamp-2 leading-relaxed flex-grow">
            {description || `${category || "Live event"} at ${venueName || city}`}
          </p>

          <div className="flex flex-wrap gap-2 mt-auto">
            {eventTags.map((tag: string, index: number) => (
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
