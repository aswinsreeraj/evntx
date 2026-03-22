import { useSearchParams } from "react-router-dom"
import { useEvents } from "../hooks"
import EventGrid from "../components/EventGrid"
import FilterSidebar from "../components/FilterSidebar"
import SortDropdown from "../components/SortDropdown"
import OrganizerCTA from "../../home/components/OrganizerCTA"
import { Search } from "lucide-react"
import { useState } from "react"

function EventListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  
  const search = searchParams.get("search") || ""
  const city = searchParams.get("city") || ""
  const category = searchParams.get("category") || ""
  const sort = searchParams.get("sort") || ""
  const start_date = searchParams.get("start_date") || ""
  const end_date = searchParams.get("end_date") || ""
  const min_price = searchParams.get("min_price") || ""
  const max_price = searchParams.get("max_price") || ""

  const citiesArray = city ? city.split(',') : []
  const categoriesArray = category ? category.split(',') : []

  const [searchInput, setSearchInput] = useState(search)

  const { data } = useEvents({ 
    limit: 12,
    city: city || undefined,
    category: category || undefined,
    search: search || undefined,
    sort: sort || undefined,
    start_date: start_date || undefined,
    end_date: end_date || undefined,
    min_price: min_price ? Number(min_price) : undefined,
    max_price: max_price ? Number(max_price) : undefined
  })
  const eventsList = data?.events || []
  const globalMinPrice = data?.price_range?.min !== undefined ? data.price_range.min : 0
  const globalMaxPrice = data?.price_range?.max !== undefined ? data.price_range.max : 100000
  // Update URL on search enter
  const handleSearchKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      const newParams = new URLSearchParams(searchParams)
      if (searchInput.trim()) {
        newParams.set("search", searchInput.trim())
      } else {
        newParams.delete("search")
      }
      setSearchParams(newParams)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-6 pt-8 pb-4">

        <div className="flex gap-4 items-center">
          <div className="flex-1 bg-white border border-gray-200 rounded-xl flex items-center p-2 shadow-sm focus-within:ring-2 focus-within:ring-blue-500/20 focus-within:border-blue-500 transition-all">
            <Search className="w-5 h-5 text-gray-400 ml-3" />
            <input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={handleSearchKeyDown}
              placeholder="Search events, city, category"
              className="flex-1 bg-transparent border-none outline-none px-3 py-1.5 text-gray-800 placeholder-gray-400 text-sm"
            />
          </div>
          <SortDropdown 
            currentSort={sort} 
            onSortChange={(newSort) => {
              const newParams = new URLSearchParams(searchParams)
              if (newSort) newParams.set("sort", newSort)
              else newParams.delete("sort")
              setSearchParams(newParams)
            }} 
          />
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-6 py-6 grid grid-cols-1 lg:grid-cols-3 gap-8">

        <div className="col-span-1 hidden lg:block">
          <FilterSidebar 
            selectedCities={citiesArray}
            selectedCategories={categoriesArray}
            startDate={start_date}
            endDate={end_date}
            minPrice={min_price}
            maxPrice={max_price}
            globalMinPrice={globalMinPrice}
            globalMaxPrice={globalMaxPrice}
            onFilterChange={(filters) => {
              const newParams = new URLSearchParams(searchParams)
              
              if (filters.cities.length > 0) newParams.set("city", filters.cities.join(","))
              else newParams.delete("city")
              
              if (filters.categories.length > 0) newParams.set("category", filters.categories.join(","))
              else newParams.delete("category")
              
              if (filters.startDate) newParams.set("start_date", filters.startDate)
              else newParams.delete("start_date")

              if (filters.endDate) newParams.set("end_date", filters.endDate)
              else newParams.delete("end_date")

              if (filters.minPrice) newParams.set("min_price", filters.minPrice)
              else newParams.delete("min_price")

              if (filters.maxPrice) newParams.set("max_price", filters.maxPrice)
              else newParams.delete("max_price")
              
              setSearchParams(newParams)
            }}
          />
        </div>


        <div className="col-span-1 lg:col-span-2">
          {eventsList.length === 0 ? (
            <div className="text-center py-20 text-gray-500">
              No events found matching your criteria.
            </div>
          ) : (
            <EventGrid events={eventsList} />
          )}
        </div>
      </div>

      <OrganizerCTA />
    </div>
  )
}

export default EventListPage
