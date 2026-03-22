import AuthModal from "../../modules/auth/components/AuthModal"
import { Bell, CircleUserRound } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import { useAuthStore } from "../../modules/auth/store/authStore"
import { Link, useNavigate } from "react-router-dom"
import { authApi } from "../../modules/auth/api"

function Navbar() {
  const [showNotifications, setShowNotifications] = useState(false)
  const [showProfileMenu, setShowProfileMenu] = useState(false)
  const { isAuthenticated, logout, roles, authModalOpen, openAuthModal, closeAuthModal } = useAuthStore()
  const navigate = useNavigate()
  const notificationRef = useRef<HTMLDivElement>(null)
  const profileMenuRef = useRef<HTMLDivElement>(null)
  const isAdmin = roles.includes("admin")
  const isOrganizer = roles.includes("organizer")
  const isGoer = isAuthenticated && !isAdmin && !isOrganizer
  const currentRole = isAdmin ? "admin" : isOrganizer ? "organizer" : isGoer ? "goer" : null
  const notifications = [
    { id: 1, title: "Booking confirmed", body: "Your ticket for Sand Castle Workshop is active." },
    { id: 2, title: "Event reminder", body: "Premium Roy by Shreya starts tomorrow at 6:00 PM." },
  ]

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (notificationRef.current && !notificationRef.current.contains(event.target as Node)) {
        setShowNotifications(false)
      }
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

              <div className="relative" ref={notificationRef}>
                <button
                  type="button"
                  onClick={() => setShowNotifications((current) => !current)}
                  className="relative rounded-full p-2 text-[#6c7480] transition hover:bg-[#f5f5f5]"
                >
                  <Bell className="h-5 w-5" />
                  <span className="absolute right-1 top-1 h-2.5 w-2.5 rounded-full bg-[#ff445d]" />
                </button>

                {showNotifications ? (
                  <div className="absolute right-0 top-12 w-80 rounded-2xl border border-gray-100 bg-white p-3 shadow-[0_16px_40px_rgba(15,23,42,0.12)]">
                    <div className="px-2 py-1 text-sm font-semibold text-[#111827]">Notifications</div>
                    <div className="mt-2 flex flex-col gap-2">
                      {notifications.map((notification) => (
                        <div key={notification.id} className="rounded-xl bg-[#f8fafc] px-3 py-3">
                          <div className="text-sm font-medium text-[#111827]">{notification.title}</div>
                          <div className="mt-1 text-xs text-[#6b7280]">{notification.body}</div>
                        </div>
                      ))}
                    </div>
                  </div>
                ) : null}
              </div>

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
