import type { ReactNode } from "react";
import { Bell, Search, User } from "lucide-react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../auth/store/authStore";
import { tokenManager } from "../../../services/tokenManager";

type Props = {
  children: ReactNode;
  title: string;
};

export default function AdminLayout({ children, title }: Props) {
  const location = useLocation();
  const navigate = useNavigate();
  const currentPath = location.pathname;
  const { logout } = useAuthStore();

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
                    isActive || (link.name === "Users" && currentPath === "/admin") // fallback active for this preview
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


            <button className="text-gray-400 hover:text-blue-500 transition-colors">
              <Bell className="w-5 h-5" fill="currentColor" />
            </button>


            <button className="w-8 h-8 rounded-full bg-gray-200 text-gray-400 flex items-center justify-center overflow-hidden">
               <User className="w-5 h-5" fill="currentColor" />
            </button>
          </div>
        </header>


        <main className="flex-1 p-8">
          {children}
        </main>
      </div>
    </div>
  );
}
