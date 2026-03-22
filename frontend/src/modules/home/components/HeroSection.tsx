import { useState, useRef } from "react";
import { Search } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function HeroSection() {
  const navigate = useNavigate();
  const [searchText, setSearchText] = useState("");
  const [showError, setShowError] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && searchText.trim()) {
      navigate(`/events?search=${encodeURIComponent(searchText.trim())}`);
    }
  };

  return (
    <div className="relative bg-black w-full h-[500px]">
      <div
        className="absolute inset-0 z-0 bg-cover bg-center"
        style={{ backgroundImage: "url('/assets/images/hero.png')" }}
      />
      <div className="absolute inset-0 bg-gradient-to-r from-purple-900/40 via-pink-800/40 to-orange-600/40 z-0" />
      <div className="absolute inset-0 bg-black/50 z-0" />

      <div className="relative z-10 flex flex-col items-center justify-center h-full text-center px-4 max-w-4xl mx-auto">
        <h1 className="text-white text-5xl md:text-6xl font-bold mb-4 tracking-tight drop-shadow-md">
          Discover Experiences <br /> That Matter
        </h1>
        <p className="text-white text-lg md:text-xl font-medium mb-10 drop-shadow">
          Concerts, tech conferences, workshops and more.
        </p>

        <div className={`bg-white rounded-full flex items-center p-2 w-full max-w-2xl shadow-lg transition-all border-2 ${showError ? "border-red-500 animate-pulse" : "border-transparent"}`}>
          <Search className="w-6 h-6 text-gray-400 ml-4" />
          <input
            ref={inputRef}
            type="text"
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value);
              if (showError) setShowError(false);
            }}
            onKeyDown={handleSearch}
            placeholder="Search events, city, category"
            className="flex-1 bg-transparent border-none outline-none px-4 py-2 text-gray-800 placeholder-gray-500 text-lg"
          />
        </div>

        <div className="flex flex-wrap justify-center gap-4 mt-8">
          <button
            onClick={() => navigate("/events")}
            className="bg-gray-900/90 border border-gray-700 backdrop-blur-sm text-white hover:bg-black px-6 py-2.5 rounded-full font-medium transition-colors">
            Browse all events
          </button>
          <button
            onClick={() => {
              if (!searchText.trim()) {
                setShowError(true);
                inputRef.current?.focus();
              } else {
                navigate(`/events?search=${encodeURIComponent(searchText.trim())}`);
              }
            }}
            className="bg-gray-900/90 border border-gray-700 backdrop-blur-sm text-white hover:bg-black px-6 py-2.5 rounded-full font-medium flex items-center gap-2 transition-colors">
            <Search className="w-4 h-4" />
            Search
          </button>
        </div>
      </div>
    </div>
  );
}