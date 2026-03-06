import EventCarousel from "./EventCarousel";

export default function FeaturedEvents() {
  // Dummy data representing featured events from Figma design
  const dummyEvents = [
    {
      id: 1,
      title: "Advancing Passive Fire Protection",
      start_time: "2026-02-19T19:00:00Z", // 07:00 PM
      city: "Bengaluru",
      cover_image_url: "https://images.unsplash.com/photo-1540575467063-178a50c2df87?w=800&q=80",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Music",
      tags: ["Music", "Party", "Dance"]
    },
    {
      id: 2,
      title: "Scorpions Coming Home Live 2026",
      start_time: "2026-02-21T12:00:00Z", // 12:00 PM
      city: "Pune",
      cover_image_url: "https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=800&q=80",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Art",
      tags: ["Art", "Workshop"]
    },
    {
      id: 3,
      title: "Anime Lovers Meet-up",
      start_time: "2026-03-03T20:00:00Z", // 08:00 PM
      city: "Chennai",
      cover_image_url: "https://images.unsplash.com/photo-1580477667995-15608129bd41?w=800&q=80",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Comedy",
      tags: ["Comedy", "Live Show"]
    }
  ];

  return <EventCarousel title="Featured Events" events={dummyEvents} />;
}
