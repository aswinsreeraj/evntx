import EventCarousel from "./EventCarousel";

export default function FeaturedEvents() {
  const dummyEvents = [
    {
      id: 1,
      title: "Advancing Passive Fire Protection",
      start_time: "2026-02-19T19:00:00Z",
      city: "Bengaluru",
      cover_image_url: "/assets/images/fire-protect.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Music",
      tags: ["Music", "Party", "Dance"]
    },
    {
      id: 2,
      title: "Scorpions Coming Home Live 2026",
      start_time: "2026-02-21T12:00:00Z",
      city: "Pune",
      cover_image_url: "/assets/images/scorpions.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Art",
      tags: ["Art", "Workshop"]
    },
    {
      id: 3,
      title: "Anime Lovers Meet-up",
      start_time: "2026-03-03T20:00:00Z", 
      city: "Chennai",
      cover_image_url: "/assets/images/anime-lover.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Comedy",
      tags: ["Comedy", "Live Show"]
    }
  ];

  return <EventCarousel title="Featured Events" events={dummyEvents} />;
}
