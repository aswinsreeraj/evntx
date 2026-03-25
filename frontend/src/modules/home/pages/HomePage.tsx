import { useState } from "react";
import HeroSection from "../components/HeroSection";
import CategorySection from "../components/CategorySection";
import FeaturedEvents from "../components/FeaturedEvents";
import OrganizerCTA from "../components/OrganizerCTA";
import CategoryEvents from "../components/CategoryEvents";
import UpcomingEvents from "../components/UpcomingEvents";

export default function HomePage() {
  const [activeCategory, setActiveCategory] = useState("All");

  return (
    <div className="min-h-screen bg-gray-50">
      <HeroSection />
      <CategorySection activeCategory={activeCategory} onCategoryChange={setActiveCategory} />
      <CategoryEvents activeCategory={activeCategory} />
      <FeaturedEvents />
      <UpcomingEvents />
      <OrganizerCTA />
    </div>
  );
}
