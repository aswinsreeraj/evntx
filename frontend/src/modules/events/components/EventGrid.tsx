import EventCard from "./EventCard"

function EventGrid({ events }: any) {
  return (
    <div className="w-full grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">
      {events.map((event: any) => (
        <EventCard key={event.id ?? event.ID ?? event.slug ?? event.Slug} event={event} />
      ))}
    </div>
  )
}

export default EventGrid
