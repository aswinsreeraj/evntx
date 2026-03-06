import Button from "../ui/Button"
import AuthModal from "../../modules/auth/components/AuthModal"
import { useState } from "react"

function Navbar() {
    const [open, setOpen] = useState(false)


  return (
    <header className="border-b bg-white">
      <div className="max-w-7xl mx-auto px-6 py-4 flex justify-between items-center">
        <h1 className="font-bold text-xl tracking-wide">EVNTX</h1>

        <div className="flex items-center gap-4">
          <button className="text-red-500 font-medium">
            + Create Event
          </button>

        <Button onClick={() => setOpen(true)}>
            Login
        </Button>

        <AuthModal open={open} onClose={() => setOpen(false)} />
        </div>
      </div>
    </header>
  )
}

export default Navbar