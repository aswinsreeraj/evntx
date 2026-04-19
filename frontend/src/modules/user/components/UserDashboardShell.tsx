import { NavLink, useNavigate } from "react-router-dom"
import { useAuthStore } from "../../auth/store/authStore"
import { tokenManager } from "../../../services/tokenManager"

type Props = {
  children: React.ReactNode
}

const navigationItems = [
  { label: "Profile", to: "/profile" },
  { label: "Bookings", to: "/profile/bookings" },
  { label: "Calendar", to: "/profile/calendar" },
  { label: "Wallet", to: "/wallet" },
]

export default function UserDashboardShell({ children }: Props) {
  const { logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    tokenManager.clearToken()
    navigate("/")
  }

  return (
    <div className="min-h-screen bg-white pt-8 pb-16">
      <div className="mx-auto flex max-w-[1120px] gap-5 px-5">
        <aside className="sticky top-[73px] flex h-[calc(100vh-73px)] w-[220px] shrink-0 flex-col rounded-[24px] border border-gray-100 bg-white p-5 shadow-[0_12px_32px_rgba(15,23,42,0.06)]">
          <div className="mt-2 flex flex-col gap-2.5">
            {navigationItems.map((item) => (
              <NavLink
                key={item.label}
                to={item.to}
                className={({ isActive }) =>
                  `rounded-2xl py-2.5 text-center text-sm font-medium transition-colors ${
                    isActive ? "bg-[#d3d4d6] text-[#111111]" : "text-[#111111] hover:bg-[#f5f5f5]"
                  }`
                }
              >
                {item.label}
              </NavLink>
            ))}
          </div>

          <button
            type="button"
            onClick={handleLogout}
            className="mt-auto rounded-2xl py-2.5 text-center text-sm font-medium text-[#ff2f3f] transition hover:bg-[#fff4f5]"
          >
            Logout
          </button>
        </aside>

        <div className="min-w-0 flex-1">{children}</div>
      </div>
    </div>
  )
}
