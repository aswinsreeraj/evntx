import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../../auth/store/authStore";
import { tokenManager } from "../../../services/tokenManager";
import NotificationMenu from "../../../shared/components/NotificationMenu";
import { CircleUserRound, Plus } from "lucide-react";
import { organizerApi, organizerWalletSummaryQueryKey } from "../api";

interface OrganizerLayoutProps {
  children: React.ReactNode;
  activeTab: string;
}

export default function OrganizerLayout({ children, activeTab }: OrganizerLayoutProps) {
  const navigate = useNavigate();
  const { logout } = useAuthStore();
  const [showProfileMenu, setShowProfileMenu] = React.useState(false);
  const profileMenuRef = React.useRef<HTMLDivElement>(null);
  const { data: walletSummary } = useQuery({
    queryKey: organizerWalletSummaryQueryKey,
    queryFn: () => organizerApi.getWalletSummary(),
  });

  React.useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(event.target as Node)) {
        setShowProfileMenu(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

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
      {}
      <div className="w-[240px] shrink-0 bg-white border-r border-gray-100 flex flex-col min-h-screen sticky top-0 py-6 px-4">
        <h1 onClick={() => navigate("/")} className="font-sigmar text-2xl tracking-wide cursor-pointer mb-8 px-2">EVNTX</h1>

        <div className="flex flex-col gap-1 mb-8">
          <p className="px-4 text-xs font-bold text-gray-400 mb-2 uppercase tracking-widest">Main</p>
          <SidebarItem label="Dashboard" />
          <SidebarItem label="My Events" route="/organizer/events" />
          <SidebarItem label="Reports" />
          <SidebarItem label="Wallet" route="/organizer/wallet" />
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

      {}
      <div className="flex-1 max-h-screen overflow-y-auto relative">
         <div className="sticky top-0 z-20 border-b border-gray-100 bg-white/95 backdrop-blur-sm">
            <div className="mx-auto flex max-w-7xl items-center justify-end gap-4 px-8 py-4">
              <button
                type="button"
                onClick={() => navigate("/organizer/events/create")}
                className="flex items-center gap-2 rounded-full bg-[#111827] px-5 py-2.5 text-sm font-medium text-white transition hover:bg-black"
              >
                <Plus className="h-4 w-4" />
                Create Event
              </button>

              <button
                type="button"
                onClick={() => navigate("/organizer/wallet")}
                className="rounded-full bg-[#f4f7fb] px-5 py-2 text-sm font-medium text-[#2a2f36] transition hover:bg-[#e9eef6]"
              >
                Wallet: ₹{new Intl.NumberFormat("en-IN", {
                  minimumFractionDigits: 2,
                  maximumFractionDigits: 2,
                }).format(walletSummary?.available_balance ?? 0)}
              </button>

              <NotificationMenu />

              <div className="relative" ref={profileMenuRef}>
                <button
                  type="button"
                  onClick={() => setShowProfileMenu((current) => !current)}
                  className="rounded-full text-[#8b9098] transition hover:text-[#111827]"
                >
                  <CircleUserRound className="h-9 w-9" />
                </button>

                {showProfileMenu && (
                  <div className="absolute right-0 top-12 w-56 rounded-2xl border border-gray-100 bg-white p-2 shadow-[0_16px_40px_rgba(15,23,42,0.12)]">
                    {[
                      { label: "View Profile", to: "/organizer/profile" },
                      { label: "My Events", to: "/organizer/events" },
                      { label: "Create Event", to: "/organizer/events/create" },
                    ].map((item) => (
                      <button
                        key={item.label}
                        type="button"
                        onClick={() => {
                          setShowProfileMenu(false);
                          navigate(item.to);
                        }}
                        className="flex w-full rounded-xl px-3 py-2.5 text-left text-sm font-medium text-[#111827] transition hover:bg-[#f8fafc]"
                      >
                        {item.label}
                      </button>
                    ))}
                    <button
                      type="button"
                      onClick={handleLogout}
                      className="flex w-full rounded-xl px-3 py-2.5 text-left text-sm font-medium text-[#e53e5d] transition hover:bg-[#fff4f5]"
                    >
                      Logout
                    </button>
                  </div>
                )}
              </div>
            </div>
         </div>
         <div className="max-w-7xl mx-auto w-full">
            {children}
         </div>
      </div>
    </div>
  );
}
