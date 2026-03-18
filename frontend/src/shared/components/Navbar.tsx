import Button from "../ui/Button"
import AuthModal from "../../modules/auth/components/AuthModal"
import { useState } from "react"
import { useAuthStore } from "../../modules/auth/store/authStore"
import { useNavigate } from "react-router-dom"

function Navbar() {
  const [authMode, setAuthMode] = useState<"goer" | "organizer" | null>(null)
  const { isAuthenticated, logout } = useAuthStore()
  const navigate = useNavigate()


  return (
    <header className="sticky top-0 z-50 bg-white shadow-sm border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-6 py-4 flex justify-between items-center">
        <h1 onClick={() => navigate("/")} className="font-sigmar text-2xl tracking-wide cursor-pointer">EVNTX</h1>

        <div className="flex items-center gap-4">
          <button className="text-red-500 font-medium" onClick={() => setAuthMode("organizer")}>
            + Create Event
          </button>

          {isAuthenticated ? (
            <Button onClick={logout}>
              Logout
            </Button>
          ) : (
            <Button onClick={() => setAuthMode("goer")}>
              Login
            </Button>
          )}

          <AuthModal 
            open={authMode !== null} 
            onClose={() => setAuthMode(null)} 
            isOrganizer={authMode === "organizer"} 
          />
        </div>
      </div>
    </header>
  )
}

export default Navbar