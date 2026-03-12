export default function CategorySection() {
  const categories = [
    { name: "Music", active: false },
    { name: "Tech", active: true },
    { name: "Business", active: false },
    { name: "Arts", active: false },
    { name: "Sports", active: false }
  ];

  return (
    <section className="max-w-7xl mx-auto px-6 mt-16">
      <h2 className="text-2xl font-bold mb-6 text-gray-900">Categories</h2>
      <div className="flex gap-4 overflow-x-auto pb-2 scrollbar-hide" style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}>
        {categories.map((category) => (
          <button
            key={category.name}
            className={`px-6 py-2 rounded-full font-medium whitespace-nowrap transition-colors ${
              category.active 
                ? "bg-[#F0BCC5] text-[#E7364D] opacity-100" 
                : "bg-gray-100 text-gray-700 hover:bg-gray-200 border border-transparent"
            }`}
          >
            {category.name}
          </button>
        ))}
      </div>
    </section>
  );
}
