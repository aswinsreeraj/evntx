import { ChevronDown } from "lucide-react";

export default function SortDropdown() {
  return (
    <div className="relative w-64">
      <button className="w-full bg-[#0b101e] text-white px-5 py-3 rounded-xl flex items-center justify-between text-sm font-medium hover:bg-black transition-colors">
        <span>Sort by</span>
        <ChevronDown className="w-4 h-4 text-white opacity-80" />
      </button>
    </div>
  );
}
