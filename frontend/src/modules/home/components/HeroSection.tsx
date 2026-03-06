function HeroSection() {
  return (
    <section
      className="h-[420px] flex items-center justify-center text-center text-white relative"
      style={{
        backgroundImage: "url('/hero.jpg')",
        backgroundSize: "cover",
        backgroundPosition: "center"
      }}
    >
      <div className="bg-black/40 absolute inset-0"></div>

      <div className="relative max-w-2xl">
        <h1 className="text-4xl font-bold mb-4">
          Discover Experiences That Matter
        </h1>

        <p className="mb-6">
          Concerts, tech conferences, workshops and more
        </p>

        <input
          placeholder="Search events, city, category"
          className="w-full p-4 rounded-xl text-black"
        />

        <div className="flex gap-4 justify-center mt-4">
          <button className="bg-black px-5 py-2 rounded-lg">
            Browse all events
          </button>

          <button className="bg-gray-800 px-5 py-2 rounded-lg">
            Browse by location
          </button>
        </div>
      </div>
    </section>
  )
}

export default HeroSection