import { ChevronDown, X, Calendar } from "lucide-react";
import { useState, useEffect } from "react";

interface FilterSidebarProps {
  selectedCities: string[];
  selectedCategories: string[];
  startDate: string;
  endDate: string;
  minPrice: string;
  maxPrice: string;
  globalMinPrice: number;
  globalMaxPrice: number;
  onFilterChange: (filters: {
    cities: string[];
    categories: string[];
    startDate: string;
    endDate: string;
    minPrice: string;
    maxPrice: string;
  }) => void;
}

export default function FilterSidebar({
  selectedCities,
  selectedCategories,
  startDate,
  endDate,
  minPrice,
  maxPrice,
  globalMinPrice,
  globalMaxPrice,
  onFilterChange,
}: FilterSidebarProps) {
  const [localCities, setLocalCities] = useState<string[]>(selectedCities || []);
  const [localCategories, setLocalCategories] = useState<string[]>(selectedCategories || []);
  const [localStartDate, setLocalStartDate] = useState(startDate || "");
  const [localEndDate, setLocalEndDate] = useState(endDate || "");

  const [localMinPrice, setLocalMinPrice] = useState(minPrice || "");
  const [localMaxPrice, setLocalMaxPrice] = useState(maxPrice || "");

  const [minVal, setMinVal] = useState(Number(minPrice) || globalMinPrice);
  const [maxVal, setMaxVal] = useState(Number(maxPrice) || globalMaxPrice);

  const [showCatDropdown, setShowCatDropdown] = useState(false);
  const [showCityDropdown, setShowCityDropdown] = useState(false);

  const categories = ["Comedy", "Music", "Workshop", "Conference"];
  const cities = ["Kochi", "Chennai", "Bangalore", "Mumbai", "Delhi", "Hyderabad"];

  useEffect(() => {
     if (!minPrice) setMinVal(globalMinPrice);
     else setMinVal(Number(minPrice));

     if (!maxPrice) setMaxVal(globalMaxPrice);
     else setMaxVal(Number(maxPrice));
  }, [minPrice, maxPrice, globalMinPrice, globalMaxPrice]);

  const handleApply = () => {
    onFilterChange({
      cities: localCities,
      categories: localCategories,
      startDate: localStartDate,
      endDate: localEndDate,
      minPrice: localMinPrice,
      maxPrice: localMaxPrice,
    });
  };

  const handleClear = () => {
    setLocalCities([]);
    setLocalCategories([]);
    setLocalStartDate("");
    setLocalEndDate("");
    setLocalMinPrice("");
    setLocalMaxPrice("");
    setMinVal(globalMinPrice);
    setMaxVal(globalMaxPrice);
    onFilterChange({
      cities: [],
      categories: [],
      startDate: "",
      endDate: "",
      minPrice: "",
      maxPrice: "",
    });
  };

  const toggleCategory = (cat: string) => {
    setLocalCategories(prev =>
      prev.includes(cat) ? prev.filter(c => c !== cat) : [...prev, cat]
    );
  };

  const toggleCity = (city: string) => {
    setLocalCities(prev =>
      prev.includes(city) ? prev.filter(c => c !== city) : [...prev, city]
    );
  };

  const handleMinChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = Math.min(Number(e.target.value), maxVal - 1);
    setMinVal(value);
    setLocalMinPrice(value.toString());
  };

  const handleMaxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = Math.max(Number(e.target.value), minVal + 1);
    setMaxVal(value);
    setLocalMaxPrice(value.toString());
  };

  const minPercent = globalMaxPrice > globalMinPrice ? ((minVal - globalMinPrice) / (globalMaxPrice - globalMinPrice)) * 100 : 0;
  const maxPercent = globalMaxPrice > globalMinPrice ? ((maxVal - globalMinPrice) / (globalMaxPrice - globalMinPrice)) * 100 : 100;

  return (
    <div className="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 flex flex-col gap-6 w-full max-w-sm">

      {}
      <div>
        <h4 className="text-sm font-medium text-gray-800 mb-3 block">Date Range</h4>
        <div className="bg-red-50 rounded-xl px-4 py-3 flex items-center gap-3 relative">
           <Calendar className="w-5 h-5 text-red-500 hidden sm:block" />
           <div className="flex gap-2 items-center flex-1 w-full relative">
              <div className="flex flex-col flex-1 relative">
                <span className="text-[10px] uppercase font-bold text-red-400 absolute -top-4 left-0">Start Date</span>
                <input
                   type="date"
                   className="bg-transparent text-sm text-gray-700 outline-none w-full !p-0 mt-2"
                   value={localStartDate}
                   onChange={(e) => setLocalStartDate(e.target.value)}
                />
              </div>
              <span className="text-gray-300 mx-1 mt-2">-</span>
              <div className="flex flex-col flex-1 relative">
                <span className="text-[10px] uppercase font-bold text-red-400 absolute -top-4 left-0">End Date</span>
                <input
                   type="date"
                   className="bg-transparent text-sm text-gray-700 outline-none w-full !p-0 mt-2"
                   value={localEndDate}
                   onChange={(e) => setLocalEndDate(e.target.value)}
                />
              </div>
           </div>
        </div>
      </div>

      {}
      <div>
        <h4 className="text-sm font-medium text-gray-800 mb-3 block">Price Range</h4>

        <div className="relative h-1.5 w-full bg-red-100 rounded-full mb-6 mt-4">
          <div
            className="absolute top-0 h-1.5 bg-[#ef4444] rounded-full"
            style={{ left: `${minPercent}%`, right: `${100 - maxPercent}%` }}
          />
          <input
            type="range"
            min={globalMinPrice}
            max={globalMaxPrice}
            value={minVal}
            onChange={handleMinChange}
            className="absolute top-[-6px] left-0 w-full appearance-none pointer-events-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-[#ef4444] [&::-webkit-slider-thumb]:border-[3px] [&::-webkit-slider-thumb]:border-white lg:cursor-pointer z-10"
          />
          <input
            type="range"
            min={globalMinPrice}
            max={globalMaxPrice}
            value={maxVal}
            onChange={handleMaxChange}
            className="absolute top-[-6px] left-0 w-full appearance-none pointer-events-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-[#ef4444] [&::-webkit-slider-thumb]:border-[3px] [&::-webkit-slider-thumb]:border-white lg:cursor-pointer z-20"
          />
        </div>

        <div className="flex items-center gap-3">
            <div className="flex-1 bg-red-50 rounded-xl px-4 py-2.5 flex items-center gap-2">
                <span className="text-red-400 font-medium tracking-tight">₹</span>
                <input
                   type="number"
                   value={localMinPrice}
                   onChange={(e) => {
                       setLocalMinPrice(e.target.value);
                       setMinVal(Number(e.target.value));
                   }}
                   className="bg-transparent outline-none w-full text-sm font-medium text-gray-800"
                   placeholder={globalMinPrice.toString()}
                />
            </div>
            <span className="text-gray-300">-</span>
            <div className="flex-1 bg-red-50 rounded-xl px-4 py-2.5 flex items-center gap-2">
                <span className="text-red-400 font-medium tracking-tight">₹</span>
                <input
                   type="number"
                   value={localMaxPrice}
                   onChange={(e) => {
                       setLocalMaxPrice(e.target.value);
                       setMaxVal(Number(e.target.value));
                   }}
                   className="bg-transparent outline-none w-full text-sm font-medium text-gray-800"
                   placeholder={globalMaxPrice.toString()}
                />
            </div>
        </div>
      </div>

      {}
      <div>
        <h4 className="text-sm font-medium text-gray-800 mb-3 block">Category</h4>
        <div className="relative">
          <button
            onClick={() => setShowCatDropdown(!showCatDropdown)}
            className="w-full bg-red-50 text-gray-500 text-sm px-4 py-3 rounded-xl flex items-center justify-between mb-3 hover:bg-red-100 transition-colors"
          >
            <div className="flex items-center gap-2 truncate">
              <svg className="w-5 h-5 text-red-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
              <span className="truncate">Select Category</span>
            </div>
            <ChevronDown className="w-5 h-5 flex-shrink-0 text-red-500" />
          </button>

          {showCatDropdown && (
            <div className="absolute z-10 w-full mt-1 bg-white border border-gray-100 rounded-xl shadow-lg max-h-48 overflow-auto">
              {categories.map(c => (
                <div
                  key={c}
                  className="px-4 py-2.5 text-sm text-gray-700 hover:bg-red-50 cursor-pointer flex justify-between items-center"
                  onClick={() => {
                    toggleCategory(c);
                  }}
                >
                  {c}
                  {localCategories.includes(c) && <X className="w-4 h-4 text-red-500" />}
                </div>
              ))}
            </div>
          )}
        </div>

        {localCategories.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {localCategories.map(cat => (
              <span key={cat} className="inline-flex items-center gap-1.5 bg-gray-100 text-gray-700 px-3 py-1.5 rounded-full text-sm">
                <X
                  className="w-4 h-4 text-red-500 cursor-pointer hover:text-red-700"
                  onClick={() => toggleCategory(cat)}
                />
                {cat}
              </span>
            ))}
          </div>
        )}
      </div>

      {}
      <div>
        <h4 className="text-sm font-medium text-gray-800 mb-3 block">Location</h4>
        <div className="relative">
          <button
            onClick={() => setShowCityDropdown(!showCityDropdown)}
            className="w-full bg-red-50 text-gray-500 text-sm px-4 py-3 rounded-xl flex items-center justify-between mb-3 hover:bg-red-100 transition-colors"
          >
            <div className="flex items-center gap-2 truncate">
              <svg className="w-5 h-5 text-red-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
              <span className="truncate">Choose locations</span>
            </div>
            <ChevronDown className="w-5 h-5 flex-shrink-0 text-red-500" />
          </button>

          {showCityDropdown && (
            <div className="absolute z-10 w-full mt-1 bg-white border border-gray-100 rounded-xl shadow-lg max-h-48 overflow-auto">
              {cities.map(c => (
                <div
                  key={c}
                  className="px-4 py-2.5 text-sm text-gray-700 hover:bg-red-50 cursor-pointer flex justify-between items-center"
                  onClick={() => {
                    toggleCity(c);
                  }}
                >
                  {c}
                  {localCities.includes(c) && <X className="w-4 h-4 text-red-500" />}
                </div>
              ))}
            </div>
          )}
        </div>

        {localCities.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {localCities.map(city => (
              <span key={city} className="inline-flex items-center gap-1.5 bg-gray-100 text-gray-700 px-3 py-1.5 rounded-full text-sm">
                <X
                  className="w-4 h-4 text-red-500 cursor-pointer hover:text-red-700"
                  onClick={() => toggleCity(city)}
                />
                {city}
              </span>
            ))}
          </div>
        )}
      </div>

      <div className="flex gap-2 mt-2">
        <button
          onClick={handleClear}
          className="flex-[1] bg-gray-100 text-gray-700 py-3.5 rounded-xl text-sm font-semibold hover:bg-gray-200 transition-colors"
        >
          Clear
        </button>
        <button
          onClick={handleApply}
          className="flex-[2] bg-[#0b101e] text-white py-3.5 rounded-xl text-sm font-semibold hover:bg-black transition-colors"
        >
          Apply Filter
        </button>
      </div>

    </div>
  );
}
