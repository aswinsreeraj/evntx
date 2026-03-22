import { ChevronLeft, ChevronRight } from "lucide-react";
import EventCard from "../../events/components/EventCard";
import { useRef } from "react";

type Props = {
  title: string;
  events: any[];
};

export default function EventCarousel({ title, events }: Props) {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: 'left' | 'right') => {
    if (scrollRef.current) {
      const { scrollLeft, clientWidth } = scrollRef.current;
      const scrollTo = direction === 'left' ? scrollLeft - clientWidth + 50 : scrollLeft + clientWidth - 50;
      scrollRef.current.scrollTo({ left: scrollTo, behavior: "smooth" });
    }
  };

  return (
    <section className="max-w-7xl mx-auto px-6 mt-16 relative">
      <h2 className="text-2xl font-bold mb-6 text-gray-900">{title}</h2>

      <div className="relative group">
        <button
          onClick={() => scroll('left')}
          className="absolute -left-5 top-1/2 -translate-y-1/2 z-10 bg-white p-2 rounded-full shadow-md text-gray-500 hover:text-gray-800 transition-opacity opacity-0 group-hover:opacity-100 hidden md:block"
        >
          <ChevronLeft className="w-6 h-6" />
        </button>

        <div
          ref={scrollRef}
          className="flex gap-6 overflow-x-auto snap-x snap-mandatory scrollbar-hide pb-4"
          style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
        >
          {events.map((event, index) => (
            <div key={event.id || event.ID || event.slug || event.Slug || index} className="min-w-[300px] md:min-w-[350px] snap-start flex-shrink-0">
              <EventCard event={event} />
            </div>
          ))}
        </div>

        <button
          onClick={() => scroll('right')}
          className="absolute -right-5 top-1/2 -translate-y-1/2 z-10 bg-white p-2 rounded-full shadow-md text-gray-500 hover:text-gray-800 transition-opacity opacity-0 group-hover:opacity-100 hidden md:block"
        >
          <ChevronRight className="w-6 h-6" />
        </button>
      </div>
    </section>
  );
}
