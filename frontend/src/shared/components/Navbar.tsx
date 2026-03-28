import AuthModal from "../../modules/auth/components/AuthModal"
import { CircleUserRound } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import { useAuthStore } from "../../modules/auth/store/authStore"
import { Link, useNavigate } from "react-router-dom"
import { authApi } from "../../modules/auth/api"
import NotificationMenu from "./NotificationMenu"

function Navbar() {
  const [showProfileMenu, setShowProfileMenu] = useState(false)
  const { isAuthenticated, logout, roles, authModalOpen, openAuthModal, closeAuthModal } = useAuthStore()
  const navigate = useNavigate()
  const profileMenuRef = useRef<HTMLDivElement>(null)
  const isAdmin = roles.includes("admin")
  const isOrganizer = roles.includes("organizer")
  const isGoer = isAuthenticated && !isAdmin && !isOrganizer
  const currentRole = isAdmin ? "admin" : isOrganizer ? "organizer" : isGoer ? "goer" : null
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(event.target as Node)) {
        setShowProfileMenu(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  const handleLogout = async () => {
    try {
      await authApi.logout()
    } catch {
      logout()
    }

    navigate(isAdmin ? "/admin/login" : "/")
  }

  const profileMenuItems =
    currentRole === "goer"
      ? [
          { label: "View Profile", to: "/profile" },
          { label: "My Bookings", to: "/profile/bookings" },
          { label: "Calendar", to: "/profile/calendar" },
        ]
      : currentRole === "organizer"
        ? [
            { label: "View Profile", to: "/organizer/profile" },
            { label: "My Events", to: "/organizer/events" },
            { label: "Create Event", to: "/organizer/events/create" },
          ]
        : currentRole === "admin"
          ? [
              { label: "Users", to: "/admin/users" },
              { label: "Organizers", to: "/admin/organizers" },
              { label: "Events", to: "/admin/events" },
            ]
          : []

  return (
    <header className="sticky top-0 z-50 bg-white shadow-sm border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-6 py-4 flex justify-between items-center">
        <h1 onClick={() => navigate("/")} className="font-sigmar text-2xl tracking-wide cursor-pointer">EVNTX</h1>

        {currentRole ? (
          <>
            <nav className="hidden md:flex items-center">
              {currentRole === "admin" ? (
                <Link to="/admin/users" className="text-gray-500 hover:text-gray-900 font-medium tracking-[0.18em] uppercase text-sm">
                  Admin
                </Link>
              ) : (
                <Link to="/events" className="text-gray-500 hover:text-gray-900 font-medium tracking-[0.18em] uppercase text-sm">
                  Explore
                </Link>
              )}
            </nav>

            <div className="flex items-center gap-4">
              {currentRole !== "admin" ? (
                <div className="rounded-full bg-[#f4f7fb] px-5 py-2 text-sm font-medium text-[#2a2f36]">
                  Wallet: ₹1,000
                </div>
              ) : null}

              {currentRole === "organizer" ? (
                <button
                  type="button"
                  onClick={() => navigate("/organizer/events/create")}
                  className="rounded-full bg-[#111827] px-5 py-2 text-sm font-medium text-white transition hover:bg-black"
                >
                  + Create Event
                </button>
              ) : null}

              <NotificationMenu />

              <div className="relative" ref={profileMenuRef}>
                <button
                  type="button"
                  onClick={(event) => {
                    event.preventDefault()
                    event.stopPropagation()
                    setShowProfileMenu((current) => !current)
                  }}
                  className="rounded-full text-[#8b9098] transition hover:text-[#111827]"
                >
                  <CircleUserRound className="h-9 w-9" />
                </button>

                {showProfileMenu ? (
                  <div className="absolute right-0 top-12 w-56 rounded-2xl border border-gray-100 bg-white p-2 shadow-[0_16px_40px_rgba(15,23,42,0.12)]">
                    {profileMenuItems.map((item) => (
                      <button
                        key={item.label}
                        type="button"
                        onClick={(event) => {
                          event.preventDefault()
                          event.stopPropagation()
                          setShowProfileMenu(false)
                          navigate(item.to)
                        }}
                        className="flex w-full rounded-xl px-3 py-2.5 text-left text-sm font-medium text-[#111827] transition hover:bg-[#f8fafc]"
                      >
                        {item.label}
                      </button>
                    ))}
                    <button
                      type="button"
                      onClick={(event) => {
                        event.preventDefault()
                        event.stopPropagation()
                        void handleLogout()
                      }}
                      className="flex w-full rounded-xl px-3 py-2.5 text-left text-sm font-medium text-[#e53e5d] transition hover:bg-[#fff4f5]"
                    >
                      Logout
                    </button>
                  </div>
                ) : null}
              </div>
            </div>
          </>
        ) : (
          <div className="flex items-center gap-4">
            <button className="text-red-500 font-medium" onClick={() => openAuthModal("organizer")}>
              + Create Event
            </button>

            {isAuthenticated ? (
              <button
                type="button"
                onClick={() => void handleLogout()}
                className="px-4 py-2 rounded-lg bg-black text-white hover:bg-gray-800 transition"
              >
                Logout
              </button>
            ) : (
              <button
                type="button"
                onClick={() => openAuthModal("goer")}
                className="px-4 py-2 rounded-lg bg-black text-white hover:bg-gray-800 transition"
              >
                Login
              </button>
            )}
          </div>
        )}

        <AuthModal
          open={authModalOpen !== null}
          onClose={closeAuthModal}
          isOrganizer={authModalOpen === "organizer"}
        />
      </div>
    </header>
  )
}

export default Navbar
