import { useState, useEffect } from "react";
import { useAuthStore } from "../../auth/store/authStore";
import { tokenManager } from "../../../services/tokenManager";
import { useNavigate } from "react-router-dom";
import { Edit2, X, CalendarDays, ChevronDown } from "lucide-react";
import { userApi } from "../api";

export default function ProfilePage() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState("Profile");
  
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [mobile, setMobile] = useState("");
  const [dob, setDob] = useState("");
  const [gender, setGender] = useState("");
  const [locations, setLocations] = useState<string[]>([]);
  
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const data = await userApi.getProfile();
        const names = (data.name || "").split(" ");
        setFirstName(names[0] || "");
        setLastName(names.slice(1).join(" ") || "");
        setEmail(data.email || "");
        setMobile(data.mobile || "");
        setDob(data.dob || "");
        setGender(data.gender || "Male");
        setLocations(data.locations || []);
      } catch (err) {
        console.error("Failed to load profile", err);
      } finally {
        setLoading(false);
      }
    };
    loadProfile();
  }, []);

  const handleLogout = () => {
    logout();
    tokenManager.clearToken();
    navigate("/");
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await userApi.updateProfile({
        name: `${firstName} ${lastName}`.trim(),
        mobile,
        dob,
        gender,
        locations,
      });
      // Optionally show a success toast here
    } catch (err) {
      console.error("Failed to save profile", err);
      alert("Failed to save profile. Please try again.");
    } finally {
      setSaving(false);
    }
  };

  const SidebarItem = ({ label }: { label: string }) => {
    const isActive = activeTab === label;
    return (
      <button
        onClick={() => setActiveTab(label)}
        className={`w-full text-center py-2.5 rounded-lg text-sm font-medium transition-colors ${
          isActive
            ? "bg-gray-200 text-gray-900"
            : "text-gray-700 hover:bg-gray-100"
        }`}
      >
        {label}
      </button>
    );
  };

  return (
    <div className="min-h-screen bg-[#f8f9fa] pt-8 pb-16">
      <div className="max-w-[1200px] mx-auto px-6 flex gap-6">
        
        {/* Sidebar */}
        <div className="w-[240px] shrink-0 bg-white rounded-3xl p-6 shadow-sm border border-gray-100 flex flex-col min-h-[600px]">
          <div className="flex flex-col gap-2 mt-4">
            <SidebarItem label="Profile" />
            <SidebarItem label="Bookings" />
            <SidebarItem label="Calendar" />
            <SidebarItem label="Wallet" />
          </div>
          <div className="mt-auto">
            <button 
              onClick={handleLogout}
              className="w-full text-center py-2.5 text-[#e53e5d] text-sm font-medium hover:bg-gray-50 rounded-lg transition-colors"
            >
              Logout
            </button>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 bg-white rounded-3xl p-10 shadow-sm border border-gray-100 min-h-[600px]">
          {loading ? (
            <div className="w-full h-full flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
          ) : (
            <>
          
          {/* Avatar Area */}
          <div className="flex flex-col items-center mb-12">
            <div className="w-32 h-32 rounded-full overflow-hidden mb-4 bg-gray-400">
               {/* Stand-in for abstract avatar image */}
               <svg viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg" className="w-full h-full text-white bg-gray-400">
                 <path d="M50 100C77.6142 100 100 77.6142 100 50C100 22.3858 77.6142 0 50 0C22.3858 0 0 22.3858 0 50C0 77.6142 22.3858 100 50 100Z" fill="currentColor"/>
                 <path d="M50 45C58.2843 45 65 38.2843 65 30C65 21.7157 58.2843 15 50 15C41.7157 15 35 21.7157 35 30C35 38.2843 41.7157 45 50 45Z" fill="white"/>
                 <path d="M25 85C25 65 35 55 50 55C65 55 75 65 75 85" stroke="white" strokeWidth="15" strokeLinecap="round"/>
               </svg>
            </div>
            <h2 className="text-xl font-medium tracking-wide text-gray-900">
              {user?.name || `${firstName} ${lastName}`}
            </h2>
          </div>

          <div className="max-w-[700px] mx-auto">
            {/* Account Details Section */}
            <div className="mb-10">
              <h3 className="text-base font-semibold text-gray-900 mb-4">Account Details</h3>
              <div className="grid grid-cols-2 gap-6">
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <label className="text-sm text-gray-600">Email</label>
                    <button className="text-[#e53e5d]"><Edit2 className="w-4 h-4" /></button>
                  </div>
                  <input 
                    type="email" 
                    readOnly
                    value={email || "randomemail@example.com"} 
                    className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-600 bg-white"
                  />
                </div>
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <label className="text-sm text-gray-600">Mobile Number</label>
                    <button className="text-[#e53e5d]"><Edit2 className="w-4 h-4" /></button>
                  </div>
                  <input 
                    type="text" 
                    value={mobile}
                    onChange={(e) => setMobile(e.target.value)}
                    className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400"
                  />
                </div>
              </div>
            </div>

            {/* Personal Details Section */}
            <div>
              <h3 className="text-base font-semibold text-gray-900 mb-4">Personal Details</h3>
              <div className="grid grid-cols-2 gap-6 mb-6">
                <div>
                  <label className="block text-sm text-gray-600 mb-2">First Name</label>
                  <input 
                    type="text" 
                    value={firstName}
                    onChange={(e) => setFirstName(e.target.value)}
                    className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400"
                  />
                </div>
                <div>
                  <label className="block text-sm text-gray-600 mb-2">Last Name</label>
                  <input 
                    type="text" 
                    value={lastName}
                    onChange={(e) => setLastName(e.target.value)}
                    className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-6 mb-8">
                <div>
                  <label className="block text-sm text-gray-600 mb-2">Date of Birth</label>
                  <div className="relative">
                    <input 
                      type="date" 
                      value={dob}
                      onChange={(e) => setDob(e.target.value)}
                      className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400 appearance-none"
                    />
                    <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none">
                      <CalendarDays className="w-5 h-5 text-[#e53e5d]" />
                    </div>
                  </div>
                </div>
                <div>
                  <label className="block text-sm text-gray-600 mb-2">Gender</label>
                  <div className="flex gap-2">
                    <button 
                      onClick={() => setGender("Male")}
                      className={`flex-1 py-3 text-sm font-medium rounded-xl border transition-colors ${gender === "Male" ? "bg-[#e53e5d] text-white border-transparent" : "border-gray-200 text-gray-600"}`}
                    >
                      Male
                    </button>
                    <button 
                      onClick={() => setGender("Female")}
                      className={`flex-1 py-3 text-sm font-medium rounded-xl border transition-colors ${gender === "Female" ? "bg-[#e53e5d] text-white border-transparent" : "border-gray-200 text-gray-600"}`}
                    >
                      Female
                    </button>
                    <div className="relative flex-1">
                      <button 
                         onClick={() => setGender("Other")}
                         className={`w-full h-full flex items-center justify-between px-3 py-3 text-sm font-medium rounded-xl border transition-colors ${gender === "Other" ? "bg-[#e53e5d] text-white border-transparent" : "border-gray-200 text-gray-600"}`}
                      >
                        Other
                        <ChevronDown className={`w-4 h-4 ${gender === "Other" ? "text-white" : "text-gray-400"}`} />
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              {/* Location Preference */}
              <div className="mb-4">
                <label className="block text-sm text-gray-600 mb-2">Preferred Location</label>
                <div className="flex gap-4 mb-2">
                  <div className="relative w-full">
                    <select 
                      onChange={(e) => {
                        const val = e.target.value;
                        if (val && !locations.includes(val)) setLocations([...locations, val]);
                        e.target.value = ""; // reset
                      }}
                      className="w-full appearance-none border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400"
                    >
                      <option value="">Select a location to add...</option>
                      <option value="Kochi">Kochi</option>
                      <option value="Bangalore">Bangalore</option>
                      <option value="Mumbai">Mumbai</option>
                    </select>
                    <ChevronDown className="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
                  </div>
                  
                  <div className="flex gap-4 w-full">
                  {/* Selected Locations Pills */}
                    {locations.map(loc => (
                      <div key={loc} className="flex-1 flex items-center justify-between border border-gray-200 rounded-xl px-4 py-3 bg-white">
                        <span className="text-sm font-medium text-gray-700">{loc}</span>
                        <button 
                          onClick={() => setLocations(locations.filter(l => l !== loc))}
                          className="hover:bg-gray-100 rounded-full p-0.5"
                        >
                          <X className="w-4 h-4 text-[#e53e5d]" />
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              
            </div>
            
          </div>
          
          <div className="max-w-[700px] mx-auto flex justify-end gap-4 mt-8">
            <button className="px-6 py-2.5 rounded-xl border border-gray-900 text-sm font-medium text-gray-900 hover:bg-gray-50 transition-colors">
              Reset Changes
            </button>
            <button 
              onClick={handleSave} 
              disabled={saving}
              className={`px-6 py-2.5 rounded-xl bg-[#0b101e] text-sm font-medium text-white hover:bg-black transition-colors ${saving ? "opacity-70 cursor-not-allowed" : ""}`}
            >
              {saving ? "Saving..." : "Apply Changes"}
            </button>
          </div>
          </>
        )}
        </div>
      </div>
    </div>
  );
}