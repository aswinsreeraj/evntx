import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../../auth/store/authStore";

export default function OrganizerCTA() {
  const navigate = useNavigate();
  const { roles, openAuthModal } = useAuthStore();
  
  const handleCreateEventClick = () => {
    if (roles && roles.some(r => r.toLowerCase() === "organizer")) {
      navigate("/organizer/events/create");
    } else {
      openAuthModal("organizer");
    }
  };
  return (
    <section className="bg-gray-100 mt-20">
      <div className="max-w-7xl mx-auto px-6 py-16 flex flex-col md:flex-row justify-between items-center gap-8">
        <div className="flex-1 md:pr-12">
          <h2 className="text-3xl md:text-4xl font-semibold mb-6 text-gray-900">
            Are you an event organizer?
          </h2>

          <div className="flex flex-wrap items-center gap-4">
            <button
              onClick={handleCreateEventClick}
              className="bg-gray-900 text-white px-6 py-3 rounded-lg font-medium hover:bg-black transition-colors"
            >
              + Create Event
            </button>
            <p className="text-gray-700 text-lg md:text-xl">
              Create. <span className="text-red-500 font-medium">Promote.</span> Sell.
            </p>
          </div>
        </div>

        <div className="flex-1 w-full flex justify-end">
          <img
            src="https://images.unsplash.com/photo-1543269865-cbf427effbad?w=800&q=80"
            alt="Event Organizer"
            className="w-full max-w-lg h-64 md:h-80 object-cover rounded-xl shadow-lg"
          />
        </div>
      </div>
    </section>
  )
}
