import type { ReactNode } from "react";
import { Search, User } from "lucide-react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../auth/store/authStore";
import { tokenManager } from "../../../services/tokenManager";
import { useEffect, useRef, useState } from "react";
import NotificationMenu from "../../../shared/components/NotificationMenu";

type Props = {
  children: ReactNode;
  title: string;
};

export default function AdminLayout({ children, title }: Props) {
  const location = useLocation();
  const navigate = useNavigate();
  const currentPath = location.pathname;
  const { logout } = useAuthStore();
  const [showProfileMenu, setShowProfileMenu] = useState(false);
  const profileMenuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(event.target as Node)) {
        setShowProfileMenu(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const platformLinks = [
    { name: "Dashboard", path: "/admin/dashboard" },
    { name: "Users", path: "/admin/users" },
    { name: "Organizers", path: "/admin/organizers" },
    { name: "Events", path: "/admin/events" },
    { name: "Reports", path: "/admin/reports" },
  ];

  const systemLinks = [
    { name: "Settings", path: "/admin/settings" },
    { name: "Audit Logs", path: "/admin/audit-logs" },
  ];

  return (
    <div className="min-h-screen bg-[#f8f9fa] flex">

      <div className="w-[240px] bg-[#0b101e] flex flex-col fixed h-full z-10 shrink-0">
        <div className="p-8 flex flex-col items-center border-b border-gray-800">
          <h1 className="text-white text-2xl font-black tracking-wider mb-1">EVNTX</h1>
          <p className="text-gray-400 text-[10px] tracking-wide uppercase">Admin Panel</p>
        </div>

        <div className="flex-1 overflow-y-auto py-6 px-4 flex flex-col gap-8 scrollbar-hide" style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}>

          <div className="flex flex-col gap-2">
            <h3 className="text-gray-500 text-[10px] font-bold tracking-widest px-4 mb-2">PLATFORM</h3>
            {platformLinks.map((link) => {
              const isActive = currentPath.includes(link.path);
              return (
                <Link
                  key={link.name}
                  to={link.path}
                  className={`px-4 py-2.5 rounded-xl text-sm font-medium transition-colors ${
                    isActive || (link.name === "Users" && currentPath === "/admin")
                      ? "bg-[#1c2438] text-white"
                      : "text-gray-400 hover:text-gray-200 hover:bg-[#141b2d]"
                  }`}
                >
                  {link.name}
                </Link>
              );
            })}
          </div>

          <div className="flex flex-col gap-2">
            <h3 className="text-gray-500 text-[10px] font-bold tracking-widest px-4 mb-2">SYSTEM</h3>
            {systemLinks.map((link) => {
              const isActive = currentPath === link.path;
              return (
                <Link
                  key={link.name}
                  to={link.path}
                  className={`px-4 py-2.5 rounded-xl text-sm font-medium transition-colors ${
                    isActive
                      ? "bg-[#1c2438] text-white"
                      : "text-gray-400 hover:text-gray-200 hover:bg-[#141b2d]"
                  }`}
                >
                  {link.name}
                </Link>
              );
            })}
          </div>

        </div>

        <div className="p-4 mt-auto">
          <button
            onClick={() => {
              logout();
              tokenManager.clearToken();
              navigate("/admin/login");
            }}
            className="w-full text-left px-4 py-2 text-[#e53e5d] text-sm font-bold hover:bg-[#141b2d] rounded-xl transition-colors"
          >
            Logout
          </button>
        </div>
      </div>

      <div className="flex-1 ml-[240px] flex flex-col min-h-screen">

        <header className="h-[72px] bg-white border-b border-gray-200 flex items-center justify-between px-8 shrink-0">
          <h2 className="text-xl font-bold text-gray-900">{title}</h2>

          <div className="flex items-center gap-6">

            <div className="flex items-center bg-transparent border border-gray-300 rounded-full px-4 py-2 w-[280px]">
              <Search className="w-4 h-4 text-blue-500 mr-2 shrink-0" />
              <input
                type="text"
                placeholder="Quick Search"
                className="bg-transparent border-none outline-none text-sm w-full placeholder-gray-400 text-gray-700"
              />
            </div>

            <NotificationMenu
              buttonClassName="relative rounded-full p-2 text-gray-400 transition-colors hover:bg-gray-100 hover:text-blue-500"
              iconClassName="w-5 h-5"
              panelClassName="absolute right-0 top-12 z-50 w-80 rounded-2xl border border-gray-100 bg-white p-3 shadow-[0_16px_40px_rgba(15,23,42,0.12)]"
            />

            <div className="relative" ref={profileMenuRef}>
              <button
                type="button"
                onClick={() => setShowProfileMenu((current) => !current)}
                className="w-8 h-8 rounded-full bg-gray-200 text-gray-400 flex items-center justify-center overflow-hidden"
              >
                 <User className="w-5 h-5" fill="currentColor" />
              </button>

              {showProfileMenu ? (
                <div className="absolute right-0 top-12 w-52 rounded-2xl border border-gray-100 bg-white p-2 shadow-[0_16px_40px_rgba(15,23,42,0.12)] z-99">
                  {[
                    { label: "Users", to: "/admin/users" },
                    { label: "Organizers", to: "/admin/organizers" },
                    { label: "Events", to: "/admin/events" },
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
                    onClick={() => {
                      logout();
                      tokenManager.clearToken();
                      navigate("/admin/login");
                    }}
                    className="flex w-full rounded-xl px-3 py-2.5 text-left text-sm font-medium text-[#e53e5d] transition hover:bg-[#fff4f5]"
                  >
                    Logout
                  </button>
                </div>
              ) : null}
            </div>
          </div>
        </header>

        <main className="flex-1 p-8">
          {children}
        </main>
      </div>
    </div>
  );
}
