import HeroSection from "../components/HeroSection";
import CategorySection from "../components/CategorySection";
import FeaturedEvents from "../components/FeaturedEvents";
import TodayEvents from "../components/TodayEvents";
import OrganizerCTA from "../components/OrganizerCTA";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      <HeroSection />
      <CategorySection />
      <FeaturedEvents />
      <TodayEvents />
      <OrganizerCTA />
    </div>
  );
}
