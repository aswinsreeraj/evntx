import HeroSection from "../components/HeroSection";
import CategorySection from "../components/CategorySection";
import FeaturedEvents from "../components/FeaturedEvents";
import TodayEvents from "../components/TodayEvents";
import OrganizerCTA from "../components/OrganizerCTA";
import CategoryEvents from "../components/CategoryEvents";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <HeroSection />
      <CategorySection />
      <CategoryEvents />
      <FeaturedEvents />
      <TodayEvents />
      <OrganizerCTA />
    </div>
  );
}
