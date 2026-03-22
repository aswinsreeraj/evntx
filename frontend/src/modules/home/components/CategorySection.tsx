interface CategorySectionProps {
  activeCategory: string;
  onCategoryChange: (cat: string) => void;
}

export default function CategorySection({ activeCategory, onCategoryChange }: CategorySectionProps) {
  const categories = ["All", "Music", "Tech", "Business", "Arts", "Sports"];

  return (
    <section className="max-w-7xl mx-auto px-6 mt-16">
      <h2 className="text-2xl font-bold mb-6 text-gray-900">Categories</h2>
      <div className="flex gap-4 overflow-x-auto pb-2 scrollbar-hide" style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}>
        {categories.map((name) => (
          <button
            key={name}
            onClick={() => onCategoryChange(name)}
            className={`px-6 py-2 rounded-full font-medium whitespace-nowrap transition-colors ${
              activeCategory === name 
                ? "bg-[#F0BCC5] text-[#E7364D] opacity-100" 
                : "bg-gray-100 text-gray-700 hover:bg-gray-200 border border-transparent"
            }`}
          >
            {name}
          </button>
        ))}
      </div>
    </section>
  );
}
