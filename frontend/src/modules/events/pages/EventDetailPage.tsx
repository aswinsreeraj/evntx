import { useParams } from "react-router-dom";
import { useEvent } from "../hooks";
import { CalendarDays, MapPin, Clock, Hourglass } from "lucide-react";
import { useState } from "react";

export default function EventDetailPage() {
  const { eventId } = useParams();
  const { data } = useEvent(eventId!);

  const [activeTab, setActiveTab] = useState("About");

  const event = {
    title: "Friday Night at Vapour Ladies Night",
    cover_image_url: "/assets/images/badass-bollywood.png",
    date: "Saturday, 21 February 2026",
    time: "07:00 PM",
    venue: "JLN Stadium, Kochi",
    duration: "5 hours 30 minutes",
    price: "5,000",
    about: {
      subtitle: "Saturday Bollywod Dhamaka",
      content: [
        "Duis placerat nisl at nisi luctus in rhoncus felis condimentum. Vivamus in augue et sem porttitor scelerisque at ac ex. Nam vel gravida lorem.",
        "Aliquam ultrices pretium odio nec hendrerit. Curabitur quis massa interdum, condimentum purus eu, bibendum felis. Proin libero ex, maximus et quam ut, volutpat condimentum tellus. Aliquam erat volutpat.",
        "Ut ipsum eros venenatis eu velit vitae landit bibendum massa.",
        "Cras id urna a quam viverra egestas sit amet et ante. In hac habitasse platea dictumst. Cras nec blandit nisi. Sed ac massa arcu."
      ]
    },
    host: {
      name: "Jane Doe",
      role: "Event Organizer",
      avatar: "/assets/images/host.jpg"
    },
    personnel: [
      {
        name: "Joe Smith",
        role: "Lead Performer",
        avatar: "/assets/images/perfomer.jpg" 
      },
      {
        name: "DJ Jazee",
        role: "Professional DJ",
        avatar: "/assets/images/dj.jpg" 
      }
    ]
  };

  
  
  const displayEvent = data || event;

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      <div className="max-w-7xl mx-auto px-6 py-10 grid grid-cols-1 lg:grid-cols-3 gap-8">
        

        <div className="lg:col-span-2 flex flex-col gap-6">
          <div className="w-full h-[400px] rounded-2xl overflow-hidden shadow-sm">
            <img 
              src={displayEvent.cover_image_url} 
              alt={displayEvent.title} 
              className="w-full h-full object-cover" 
            />
          </div>

          <h1 className="text-2xl md:text-3xl font-bold text-gray-900 mt-2">
            {displayEvent.title}
          </h1>

          <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mt-2">
            <div className="flex gap-8 border-b border-gray-100 mb-6">
              {["About", "Venue", "Terms & Conditions"].map((tab) => (
                <button
                  key={tab}
                  className={`pb-3 text-sm font-medium transition-colors relative ${
                    activeTab === tab ? "text-[#e53e5d]" : "text-gray-600 hover:text-gray-900"
                  }`}
                  onClick={() => setActiveTab(tab)}
                >
                  {tab}
                  {activeTab === tab && (
                    <div className="absolute bottom-0 left-0 w-full h-0.5 bg-[#e53e5d] rounded-t-full"></div>
                  )}
                </button>
              ))}
            </div>

            {activeTab === "About" && (
              <div className="animate-in fade-in duration-300">
                <h3 className="font-bold text-gray-900 text-sm mb-4">{displayEvent.about?.subtitle || displayEvent.title}</h3>
                <div className="flex flex-col gap-4 text-sm text-gray-700 leading-relaxed">
                  {displayEvent.about?.content ? displayEvent.about.content.map((paragraph: string, index: number) => (
                    <p key={index}>{paragraph}</p>
                  )) : (
                    <p>{displayEvent.description}</p>
                  )}
                </div>
              </div>
            )}
            {activeTab === "Venue" && (
              <div className="animate-in fade-in duration-300 text-sm text-gray-700">
                Venue details will go here.
              </div>
            )}
            {activeTab === "Terms & Conditions" && (
              <div className="animate-in fade-in duration-300 text-sm text-gray-700">
                Terms and conditions details will go here.
              </div>
            )}
          </div>
        </div>


        <div className="lg:col-span-1 flex flex-col gap-6">
          

          <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-6">Book Tickets</h3>
            
            <div className="flex flex-col gap-4 mb-8">
              <div className="flex items-start gap-4 text-sm text-gray-700">
                <CalendarDays className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                <span>{displayEvent.date || new Date(displayEvent.start_time).toLocaleDateString()}</span>
              </div>
              <div className="flex items-start gap-4 text-sm text-gray-700">
                <Clock className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                <span>{displayEvent.time || new Date(displayEvent.start_time).toLocaleTimeString()}</span>
              </div>
              <div className="flex items-start gap-4 text-sm text-gray-700">
                <MapPin className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                <span>{displayEvent.venue || displayEvent.city}</span>
              </div>
              <div className="flex items-start gap-4 text-sm text-gray-700">
                <Hourglass className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                <span>{displayEvent.duration || "N/A"}</span>
              </div>
            </div>

            <div className="bg-[#fcf3f4] rounded-xl p-4 mb-4 flex justify-between items-center">
              <span className="text-sm font-medium text-gray-700">Price From</span>
              <span className="text-sm font-bold text-[#e53e5d]">₹ {displayEvent.price || "N/A"}</span>
            </div>

            <button className="w-full bg-[#0b101e] hover:bg-black text-white py-3.5 rounded-xl text-sm font-medium transition-colors">
              Continue to Booking
            </button>
          </div>



          {displayEvent.host && (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-6">Host</h3>
              <div className="flex items-center gap-4 mb-6">
                <img src={displayEvent.host.avatar} alt={displayEvent.host.name} className="w-12 h-12 rounded-full object-cover" />
                <div>
                  <h4 className="font-bold text-gray-900 text-sm">{displayEvent.host.name}</h4>
                  <p className="text-xs text-gray-500">{displayEvent.host.role}</p>
                </div>
              </div>
              <button className="w-full border border-gray-300 hover:bg-gray-50 text-gray-700 py-2.5 rounded-xl text-sm font-medium transition-colors">
                View Profile
              </button>
            </div>
          )}


          {displayEvent.personnel && (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-6">Key Personnel</h3>
              <div className="flex flex-col gap-6">
                {displayEvent.personnel.map((person: any, index: number) => (
                  <div key={index} className="flex items-center gap-4">
                    <img src={person.avatar} alt={person.name} className="w-12 h-12 rounded-full object-cover" />
                    <div>
                      <h4 className="font-bold text-gray-900 text-sm">{person.name}</h4>
                      <p className="text-xs text-gray-500">{person.role}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

        </div>
      </div>
    </div>
  );
}