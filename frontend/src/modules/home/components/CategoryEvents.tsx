import EventCarousel from "./EventCarousel";

export default function CategoryEvents() {
  const dummyEvents = [
    {
      id: 1,
      title: "Saturday Bollywod Dhamaka",
      start_time: "2026-02-19T19:00:00Z",
      city: "Bengaluru",
      cover_image_url: "/assets/images/badass-bollywood.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Music",
      tags: ["Music", "Party", "Dance"]
    },
    {
      id: 2,
      title: "Sand Castle Workshop",
      start_time: "2026-02-21T12:00:00Z",
      city: "Pune",
      cover_image_url: "/assets/images/sand-castle.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Art",
      tags: ["Art", "Workshop"]
    },
    {
      id: 3,
      title: "Premium Roy by Shreya",
      start_time: "2026-03-03T20:00:00Z",
      city: "Chennai",
      cover_image_url: "/assets/images/premium-roy.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Comedy",
      tags: ["Comedy", "Live Show"]
    }
  ];

  return <EventCarousel title="" events={dummyEvents} />;
}
