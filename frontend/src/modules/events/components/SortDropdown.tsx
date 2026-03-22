import { ChevronDown, Check } from "lucide-react";
import { useState } from "react";

interface SortDropdownProps {
  currentSort: string;
  onSortChange: (sort: string) => void;
}

export default function SortDropdown({ currentSort, onSortChange }: SortDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);

  const sortOptions = [
    { label: "Date Ascending", value: "date_asc" },
    { label: "Date Descending", value: "date_desc" },
    { label: "Recently Added", value: "created_at_desc" },
  ];

  const currentOption = sortOptions.find(opt => opt.value === currentSort) || { label: "Sort by" };

  return (
    <div className="relative w-64">
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className="w-full bg-[#0b101e] text-white px-5 py-3 rounded-xl flex items-center justify-between text-sm font-medium hover:bg-black transition-colors"
      >
        <span>{currentOption.label}</span>
        <ChevronDown className={`w-4 h-4 text-white opacity-80 transition-transform ${isOpen ? "rotate-180" : ""}`} />
      </button>

      {isOpen && (
        <div className="absolute z-20 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-lg overflow-hidden">
          {sortOptions.map((option) => (
            <button
              key={option.value}
              onClick={() => {
                onSortChange(option.value);
                setIsOpen(false);
              }}
              className="w-full text-left px-5 py-3 text-sm flex items-center justify-between hover:bg-gray-50 transition-colors"
            >
              <span className={currentSort === option.value ? "font-medium text-gray-900" : "text-gray-600"}>
                {option.label}
              </span>
              {currentSort === option.value && <Check className="w-4 h-4 text-[#e53e5d]" />}
            </button>
          ))}
          {currentSort && (
            <button
              onClick={() => {
                onSortChange("");
                setIsOpen(false);
              }}
              className="w-full text-left px-5 py-3 text-sm text-red-500 hover:bg-red-50 transition-colors border-t border-gray-50"
            >
              Clear Sort
            </button>
          )}
        </div>
      )}
    </div>
  );
}

