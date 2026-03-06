function OrganizerCTA() {
  return (
    <section className="bg-gray-100 mt-20">
      <div className="max-w-7xl mx-auto px-6 py-16 flex justify-between items-center">

        <div>
          <h2 className="text-2xl font-semibold mb-3">
            Are you an event organizer?
          </h2>

          <p className="text-gray-500 mb-4">
            Create. Promote. Sell.
          </p>

          <button className="bg-black text-white px-5 py-2 rounded-lg">
            + Create Event
          </button>
        </div>

        <img
          src="/organizer.jpg"
          className="h-60 rounded-xl"
        />
      </div>
    </section>
  )
}

export default OrganizerCTA