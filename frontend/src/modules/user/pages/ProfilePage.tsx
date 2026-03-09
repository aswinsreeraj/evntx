import { useState, useEffect, useRef } from "react";
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

  const [profileImage, setProfileImage] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [uploadingImage, setUploadingImage] = useState(false);

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [successMsg, setSuccessMsg] = useState("");

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
        if (data.profile_image) {
          setProfileImage(`${import.meta.env.VITE_API_BASE_URL}${data.profile_image}`);
        }
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

  const validate = () => {
    const newErrors: Record<string, string> = {};
    if (!firstName.trim() || firstName.trim().length < 2) {
      newErrors.firstName = "First name must be at least 2 characters";
    }
    if (lastName.trim() && lastName.trim().length < 2) {
      newErrors.lastName = "Last name must be at least 2 characters";
    }
    if (mobile && !/^\+?[0-9\s-]{10,15}$/.test(mobile)) {
      newErrors.mobile = "Please enter a valid mobile number";
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSave = async () => {
    setSuccessMsg("");
    if (!validate()) return;

    setSaving(true);
    setErrors({});
    try {
      await userApi.updateProfile({
        name: `${firstName} ${lastName}`.trim(),
        mobile,
        dob,
        gender,
        locations,
      });
      setSuccessMsg("Profile updated successfully!");
      setTimeout(() => setSuccessMsg(""), 3000);
    } catch (err: any) {
      console.error("Failed to save profile", err);
      setErrors({ api: err.response?.data?.message || "Failed to save profile. Please try again." });
    } finally {
      setSaving(false);
    }
  };

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadingImage(true);
    setErrors({ ...errors, image: "" });
    try {
      const data = await userApi.uploadProfileImage(file);
      setProfileImage(`${import.meta.env.VITE_API_BASE_URL}${data.profile_image}`);
      setSuccessMsg("Profile image updated successfully!");
      setTimeout(() => setSuccessMsg(""), 3000);
    } catch (err: any) {
      console.error("Failed to upload image", err);
      setErrors({ ...errors, image: err.response?.data?.message || "Failed to upload image." });
    } finally {
      setUploadingImage(false);
    }
  };

  const SidebarItem = ({ label }: { label: string }) => {
    const isActive = activeTab === label;
    return (
      <button
        onClick={() => setActiveTab(label)}
        className={`w-full text-center py-2.5 rounded-lg text-sm font-medium transition-colors ${isActive
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


        <div className="flex-1 bg-white rounded-3xl p-10 shadow-sm border border-gray-100 min-h-[600px]">
          {loading ? (
            <div className="w-full h-full flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
          ) : (
            <>


              <div className="flex flex-col items-center mb-12">
                <div className="relative group">
                  <input
                    type="file"
                    ref={fileInputRef}
                    onChange={handleImageUpload}
                    accept="image/jpeg,image/png,image/webp"
                    className="hidden"
                  />
                  <div className="w-32 h-32 rounded-full overflow-hidden mb-4 bg-gray-200 shadow-inner relative flex items-center justify-center text-gray-400 font-bold text-3xl">
                    {profileImage ? (
                      <img src={profileImage} alt="Profile" className="w-full h-full object-cover" />
                    ) : (
                      <span>{firstName.charAt(0)}{lastName.charAt(0)}</span>
                    )}
                  </div>
                  <button
                    onClick={() => fileInputRef.current?.click()}
                    disabled={uploadingImage}
                    className="absolute right-0 bottom-4 w-9 h-9 bg-white rounded-full shadow-md border border-gray-100 flex items-center justify-center text-gray-600 hover:text-gray-900 transition-colors z-10 hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                  >
                    {uploadingImage ? (
                      <div className="w-4 h-4 border-2 border-gray-300 border-t-gray-800 rounded-full animate-spin" />
                    ) : (
                      <Edit2 className="w-4 h-4" />
                    )}
                  </button>
                </div>
                {errors.image && <p className="text-red-500 text-xs mb-2">{errors.image}</p>}
                <h2 className="text-xl font-medium tracking-wide text-gray-900">
                  {user?.name || `${firstName} ${lastName}`}
                </h2>
              </div>

              <div className="max-w-[700px] mx-auto">

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
                        onChange={(e) => { setMobile(e.target.value); setErrors({ ...errors, mobile: "" }) }}
                        className={`w-full border ${errors.mobile ? 'border-red-500' : 'border-gray-200'} rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:ring-1 ${errors.mobile ? 'focus:ring-red-400 focus:border-red-500' : 'focus:ring-gray-400 focus:border-gray-400'}`}
                      />
                      {errors.mobile && <span className="text-red-500 text-xs mt-1 block">{errors.mobile}</span>}
                    </div>
                  </div>
                </div>


                <div>
                  <h3 className="text-base font-semibold text-gray-900 mb-4">Personal Details</h3>
                  <div className="grid grid-cols-2 gap-6 mb-6">
                    <div>
                      <label className="block text-sm text-gray-600 mb-2">First Name</label>
                      <input
                        type="text"
                        value={firstName}
                        onChange={(e) => { setFirstName(e.target.value); setErrors({ ...errors, firstName: "" }) }}
                        className={`w-full border ${errors.firstName ? 'border-red-500' : 'border-gray-200'} rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:ring-1 ${errors.firstName ? 'focus:ring-red-400 focus:border-red-500' : 'focus:ring-gray-400 focus:border-gray-400'}`}
                      />
                      {errors.firstName && <span className="text-red-500 text-xs mt-1 block">{errors.firstName}</span>}
                    </div>
                    <div>
                      <label className="block text-sm text-gray-600 mb-2">Last Name</label>
                      <input
                        type="text"
                        value={lastName}
                        onChange={(e) => { setLastName(e.target.value); setErrors({ ...errors, lastName: "" }) }}
                        className={`w-full border ${errors.lastName ? 'border-red-500' : 'border-gray-200'} rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:ring-1 ${errors.lastName ? 'focus:ring-red-400 focus:border-red-500' : 'focus:ring-gray-400 focus:border-gray-400'}`}
                      />
                      {errors.lastName && <span className="text-red-500 text-xs mt-1 block">{errors.lastName}</span>}
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


                  <div className="mb-4">
                    <label className="block text-sm text-gray-600 mb-2">Preferred Location</label>
                    <div className="flex gap-4 mb-2">
                      <div className="relative w-full">
                        <select
                          onChange={(e) => {
                            const val = e.target.value;
                            if (val && !locations.includes(val)) setLocations([...locations, val]);
                            e.target.value = "";
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

              <div className="max-w-[700px] mx-auto mt-6">
                {errors.api && (
                  <div className="w-full mb-4 p-3 bg-red-50 border border-red-100 rounded-lg">
                    <p className="text-red-600 text-sm text-center font-medium">{errors.api}</p>
                  </div>
                )}
                {successMsg && (
                  <div className="w-full mb-4 p-3 bg-green-50 border border-green-100 rounded-lg">
                    <p className="text-green-600 text-sm text-center font-medium">{successMsg}</p>
                  </div>
                )}
                <div className="flex justify-end gap-4">
                  <button
                    className="px-6 py-2.5 rounded-xl border border-gray-900 text-sm font-medium text-gray-900 hover:bg-gray-50 transition-colors"
                    onClick={() => window.location.reload()}
                  >
                    Reset Changes
                  </button>
                  <button
                    onClick={handleSave}
                    disabled={saving}
                    className={`px-6 py-2.5 flex justify-center items-center rounded-xl bg-[#0b101e] text-sm font-medium text-white hover:bg-black transition-colors ${saving ? "opacity-70 cursor-not-allowed min-w-[130px]" : "min-w-[130px]"}`}
                  >
                    {saving ? (
                      <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
                    ) : (
                      "Apply Changes"
                    )}
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}