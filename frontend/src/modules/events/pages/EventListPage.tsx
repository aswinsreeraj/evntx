import { useEvents } from "../hooks"
import EventGrid from "../components/EventGrid"
import FilterSidebar from "../components/FilterSidebar"
import SortDropdown from "../components/SortDropdown"
import OrganizerCTA from "../../home/components/OrganizerCTA"
import { Search } from "lucide-react"

function EventListPage() {
  const { data, isLoading } = useEvents()

  const dummyEvents = [
    {
      id: 1,
      title: "Sand Castle Workshop",
      start_time: "2026-02-21T12:00:00Z",
      city: "Pune",
      cover_image_url: "/assets/images/sand-castle.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      tags: ["Art", "Workshop"]
    },
    {
      id: 2,
      title: "Premium Roy by Shreya",
      start_time: "2026-03-03T18:00:00Z",
      city: "Chennai",
      cover_image_url: "/assets/images/premium-roy.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      tags: ["Comedy", "Live Show"]
    },
    {
      id: 3,
      title: "Advancing Passive Fire Protection",
      start_time: "2026-02-19T19:00:00Z", 
      city: "Bengaluru",
      cover_image_url: "/assets/images/fire-protect.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      tags: ["Music", "Party", "Dance"]
    },
    {
      id: 4,
      title: "Scorpions Coming Home Live 2026",
      start_time: "2026-02-21T12:00:00Z", 
      city: "Pune",
      cover_image_url: "/assets/images/scorpions.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      tags: ["Art", "Workshop"]
    },
    {
      id: 5,
      title: "Anime Lovers Meet-up",
      start_time: "2026-03-03T20:00:00Z", 
      city: "Chennai",
      cover_image_url: "/assets/images/anime-lover.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      tags: ["Comedy", "Live Show"]
    },
    {
      id: 6,
      title: "Saturday Bollywood Dhamaka",
      start_time: "2026-02-19T19:00:00Z", 
      city: "Bengaluru",
      cover_image_url: "/assets/images/badass-bollywood.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      tags: ["Music", "Party", "Dance"]
    }
  ];

  const eventsList = dummyEvents; 

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-6 pt-8 pb-4">
        {/* Search and Sort Top Bar */}
        <div className="flex gap-4 items-center">
          <div className="flex-1 bg-white border border-gray-200 rounded-xl flex items-center p-2 shadow-sm focus-within:ring-2 focus-within:ring-blue-500/20 focus-within:border-blue-500 transition-all">
            <Search className="w-5 h-5 text-gray-400 ml-3" />
            <input
              type="text"
              placeholder="Search events, city, category"
              className="flex-1 bg-transparent border-none outline-none px-3 py-1.5 text-gray-800 placeholder-gray-400 text-sm"
            />
          </div>
          <SortDropdown />
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-6 py-6 grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Sidebar */}
        <div className="col-span-1 hidden lg:block">
          <FilterSidebar />
        </div>

        {/* Event list */}
        <div className="col-span-1 lg:col-span-2">
          <EventGrid events={eventsList} />
        </div>
      </div>

      <OrganizerCTA />
    </div>
  )
}

export default EventListPage