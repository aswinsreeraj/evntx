import React from "react";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../../auth/store/authStore";
import { tokenManager } from "../../../services/tokenManager";

interface OrganizerLayoutProps {
  children: React.ReactNode;
  activeTab: string;
}

export default function OrganizerLayout({ children, activeTab }: OrganizerLayoutProps) {
  const navigate = useNavigate();
  const { logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    tokenManager.clearToken();
    navigate("/");
  };

  const SidebarItem = ({ label, route }: { label: string; route?: string }) => {
    const isActive = activeTab === label;
    return (
      <button
        onClick={() => {
          if (route && !isActive) navigate(route);
        }}
        className={`w-full text-left px-4 py-2.5 rounded-lg text-sm font-medium transition-colors ${isActive
          ? "bg-gray-200 text-gray-900"
          : "text-gray-700 hover:bg-gray-100"
          }`}
      >
        {label}
      </button>
    );
  };

  return (
    <div className="min-h-screen border-t border-gray-100 bg-[#f8f9fa] flex">
      {/* Sidebar */}
      <div className="w-[240px] shrink-0 bg-white border-r border-gray-100 flex flex-col min-h-screen sticky top-0 py-6 px-4">
        <h1 onClick={() => navigate("/")} className="font-sigmar text-2xl tracking-wide cursor-pointer mb-8 px-2">EVNTX</h1>
        
        <div className="flex flex-col gap-1 mb-8">
          <p className="px-4 text-xs font-bold text-gray-400 mb-2 uppercase tracking-widest">Main</p>
          <SidebarItem label="Dashboard" />
          <SidebarItem label="My Events" route="/organizer/events" />
          <SidebarItem label="Reports" />
          <SidebarItem label="Wallet" />
        </div>

        <div className="flex flex-col gap-1">
          <p className="px-4 text-xs font-bold text-gray-400 mb-2 uppercase tracking-widest">Account</p>
          <SidebarItem label="Profile" route="/organizer/profile" />
        </div>
        
        <div className="mt-auto">
          <button
            onClick={handleLogout}
            className="w-full text-left px-4 py-2.5 text-[#e53e5d] text-sm font-medium hover:bg-gray-50 rounded-lg transition-colors"
          >
            Logout
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 max-h-screen overflow-y-auto relative">
         <div className="max-w-7xl mx-auto w-full">
            {children}
         </div>
      </div>
    </div>
  );
}
