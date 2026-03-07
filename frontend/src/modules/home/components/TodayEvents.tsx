import EventCarousel from "./EventCarousel";

export default function TodayEvents() {
  // Dummy data representing today's events from Figma design
  const dummyEvents = [
    {
      id: 4,
      title: "UNO Party Meetup",
      start_time: "2026-02-19T19:00:00Z",
      city: "Bengaluru",
      cover_image_url: "https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=800&q=80",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Meetup",
      tags: ["Meetup", "Dating"]
    },
    {
      id: 5,
      title: "If I'm Not Wrong By Tarang Hardikar",
      start_time: "2026-02-19T12:00:00Z",
      city: "Pune",
      cover_image_url: "https://images.unsplash.com/photo-1585699324551-f6c309eedeca?w=800&q=80",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Comedy",
      tags: ["Comedy", "Live Show"]
    },
    {
      id: 6,
      title: "Scratch Pottery",
      start_time: "2026-02-19T12:00:00Z",
      city: "Pune",
      cover_image_url: "https://images.unsplash.com/photo-1610701596007-11502861dcfa?w=800&q=80",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Art",
      tags: ["Art", "Workshop"]
    }
  ];

  return <EventCarousel title="Today's Events" events={dummyEvents} />;
}
