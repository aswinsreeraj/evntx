import { Instagram, Facebook, Twitter } from "lucide-react";

export default function Footer() {
  return (
    <footer className="bg-[#0b101e] text-gray-300">
      <div className="max-w-7xl mx-auto px-6 py-12 flex flex-col items-center">
        <h2 className="font-bold text-white text-2xl tracking-wide mb-4">EVNTX</h2>

        <div className="flex gap-4 mb-8">
          <a href="#" className="text-gray-400 hover:text-white transition-colors">
            <Instagram className="w-5 h-5" />
          </a>
          <a href="#" className="text-gray-400 hover:text-white transition-colors">
            <Facebook className="w-5 h-5" />
          </a>
          <a href="#" className="text-gray-400 hover:text-white transition-colors">
            <Twitter className="w-5 h-5" />
          </a>
        </div>

        <div className="flex flex-wrap justify-center gap-x-8 gap-y-4 text-sm mb-12 text-gray-400">
          <a href="#" className="hover:text-white transition-colors">About</a>
          <a href="#" className="hover:text-white transition-colors">Careers</a>
          <a href="#" className="hover:text-white transition-colors">Support</a>
          <a href="#" className="hover:text-white transition-colors">Terms of Service</a>
          <a href="#" className="hover:text-white transition-colors">Privacy</a>
          <a href="#" className="hover:text-white transition-colors">Contact Us</a>
        </div>

        <p className="text-xs text-gray-500">© EVNTX, 2026</p>
      </div>
    </footer>
  )
}
