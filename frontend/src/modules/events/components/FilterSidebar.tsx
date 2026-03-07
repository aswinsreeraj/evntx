import { CalendarDays, ChevronDown, X } from "lucide-react";

export default function FilterSidebar() {
  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex flex-col gap-8">
      
      {/* Date Range */}
      <div>
        <h4 className="text-sm font-medium text-gray-700 mb-3">Date Range</h4>
        <div className="relative">
          <input 
            type="text" 
            placeholder="Select event dates" 
            className="w-full bg-red-50 text-red-500 placeholder-red-300 text-sm px-10 py-2.5 rounded-lg border-none focus:outline-none focus:ring-1 focus:ring-red-200"
          />
          <CalendarDays className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-red-400" />
        </div>
      </div>

      {/* Price Range */}
      <div>
        <h4 className="text-sm font-medium text-gray-700 mb-3">Price Range</h4>
        {/* Visual Slider Placeholder */}
        <div className="h-1 bg-gray-200 w-full rounded-full mb-4 relative">
          <div className="absolute left-[10%] right-[30%] h-full bg-[#e53e5d] rounded-full"></div>
          <div className="absolute left-[10%] top-1/2 -translate-y-1/2 w-3 h-3 bg-[#e53e5d] rounded-full shadow border-2 border-white"></div>
          <div className="absolute right-[30%] top-1/2 -translate-y-1/2 w-3 h-3 bg-[#e53e5d] rounded-full shadow border-2 border-white"></div>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex-1 bg-red-50 text-red-500 rounded-lg px-3 py-2 text-sm flex items-center justify-between">
            <span className="text-red-400">₹</span>
            <span>100</span>
          </div>
          <span className="text-gray-400">-</span>
          <div className="flex-1 bg-red-50 text-red-500 rounded-lg px-3 py-2 text-sm flex items-center justify-between">
            <span className="text-red-400">₹</span>
            <span>99999</span>
          </div>
        </div>
      </div>

      {/* Category */}
      <div>
        <h4 className="text-sm font-medium text-gray-700 mb-3">Category</h4>
        <button className="w-full bg-red-50 text-red-500 text-sm px-4 py-2.5 rounded-lg flex items-center justify-between mb-3 hover:bg-red-100 transition-colors">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 border-2 border-red-400 rounded-sm"></div>
            <span>Select Category</span>
          </div>
          <ChevronDown className="w-4 h-4 text-red-400" />
        </button>
        <div className="flex flex-wrap gap-2">
          <span className="inline-flex items-center gap-1 bg-gray-100/80 text-gray-700 px-3 py-1.5 rounded-full text-xs font-medium">
            <X className="w-3 h-3 text-red-500 cursor-pointer hover:text-red-700" /> Music
          </span>
          <span className="inline-flex items-center gap-1 bg-gray-100/80 text-gray-700 px-3 py-1.5 rounded-full text-xs font-medium">
            <X className="w-3 h-3 text-red-500 cursor-pointer hover:text-red-700" /> Comedy
          </span>
        </div>
      </div>

      {/* Location */}
      <div>
        <h4 className="text-sm font-medium text-gray-700 mb-3">Location</h4>
        <button className="w-full bg-red-50 text-red-500 text-sm px-4 py-2.5 rounded-lg flex items-center justify-between mb-3 hover:bg-red-100 transition-colors">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 border-2 border-red-400 rounded-sm"></div>
            <span>Choose locations</span>
          </div>
          <ChevronDown className="w-4 h-4 text-red-400" />
        </button>
        <div className="flex flex-wrap gap-2">
          <span className="inline-flex items-center gap-1 bg-gray-100/80 text-gray-700 px-3 py-1.5 rounded-full text-xs font-medium">
            <X className="w-3 h-3 text-red-500 cursor-pointer hover:text-red-700" /> Kochi
          </span>
          <span className="inline-flex items-center gap-1 bg-gray-100/80 text-gray-700 px-3 py-1.5 rounded-full text-xs font-medium">
            <X className="w-3 h-3 text-red-500 cursor-pointer hover:text-red-700" /> Chennai
          </span>
        </div>
      </div>

      {/* Apply Filter Button */}
      <button className="w-full bg-[#0b101e] text-white py-3 rounded-xl text-sm font-medium hover:bg-black transition-colors mt-2">
        Apply Filter
      </button>

    </div>
  );
}
