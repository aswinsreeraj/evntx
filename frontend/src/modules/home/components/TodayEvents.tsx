import EventCarousel from "./EventCarousel";

export default function TodayEvents() {
  const dummyEvents = [
    {
      id: 4,
      title: "UNO Party Meetup",
      start_time: "2026-02-19T19:00:00Z",
      city: "Bengaluru",
      cover_image_url: "/assets/images/uno-party.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Meetup",
      tags: ["Meetup", "Dating"]
    },
    {
      id: 5,
      title: "If I'm Not Wrong By Tarang Hardikar",
      start_time: "2026-02-19T12:00:00Z",
      city: "Pune",
      cover_image_url: "/assets/images/if-im-not-wrong.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Comedy",
      tags: ["Comedy", "Live Show"]
    },
    {
      id: 6,
      title: "Scratch Pottery",
      start_time: "2026-02-19T12:00:00Z",
      city: "Pune",
      cover_image_url: "/assets/images/scratch-pottery.png",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus aliquam volutpat lorem, sed vehicula mi pellentesque.",
      category: "Art",
      tags: ["Art", "Workshop"]
    }
  ];

  return <EventCarousel title="Today's Events" events={dummyEvents} />;
}
