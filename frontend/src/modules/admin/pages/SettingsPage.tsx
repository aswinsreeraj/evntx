import { useState, useEffect } from "react";
import AdminLayout from "../components/AdminLayout";
import { useSettings, useUpdateSettings, usePaymentSettings, useUpdatePaymentProvider, useAdmins } from "../hooks";
import { Loader2, ShieldCheck } from "lucide-react";

export default function SettingsPage() {
  const { data: settings, isLoading: isSettingsLoading } = useSettings();
  const { data: paymentSettingsList, isLoading: isPaymentsLoading } = usePaymentSettings();
  const { data: adminsData, isLoading: isAdminsLoading } = useAdmins();

  const { mutate: updateSettings, isPending: isSaving } = useUpdateSettings();
  const { mutate: updatePayment } = useUpdatePaymentProvider();

  const isLoading = isSettingsLoading || isPaymentsLoading || isAdminsLoading;

  const [localSettings, setLocalSettings] = useState<any>(null);
  const [razorpayEnabled, setRazorpayEnabled] = useState(false);
  const [walletEnabled, setWalletEnabled] = useState(false);
  const [toastMessage, setToastMessage] = useState<{type: "success" | "error", text: string} | null>(null);

  useEffect(() => {
    if (settings) {
      setLocalSettings({ ...settings });
    }
  }, [settings]);

  useEffect(() => {
    if (paymentSettingsList) {
      const rp = paymentSettingsList.find((p) => p.provider === "razorpay");
      if (rp) {
        setRazorpayEnabled(rp.is_enabled);
      }
      const wt = paymentSettingsList.find((p) => p.provider === "wallet");
      if (wt) {
        setWalletEnabled(wt.is_enabled);
      }
    }
  }, [paymentSettingsList]);

  if (isLoading || !localSettings) {
    return (
      <AdminLayout title="Manage Platform">
        <div className="flex h-64 items-center justify-center">
          <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
        </div>
      </AdminLayout>
    );
  }

  const handleToggle = (key: string) => {
    setLocalSettings((prev: any) => ({
      ...prev,
      [key]: !prev[key],
    }));
  };

  const handleChange = (key: string, value: any) => {
    setLocalSettings((prev: any) => ({
      ...prev,
      [key]: value,
    }));
  };

  const handleSave = () => {
    updateSettings(localSettings, {
      onSuccess: () => {
        
        const rp = paymentSettingsList?.find((p) => p.provider === "razorpay");
        if (rp && rp.is_enabled !== razorpayEnabled) {
          updatePayment({
            provider: "razorpay",
            data: { is_enabled: razorpayEnabled, config: rp.config },
          });
        }
        const wt = paymentSettingsList?.find((p) => p.provider === "wallet");
        if (wt && wt.is_enabled !== walletEnabled) {
          updatePayment({
            provider: "wallet",
            data: { is_enabled: walletEnabled, config: wt.config },
          });
        }
        setToastMessage({ type: "success", text: "Settings saved successfully!" });
        setTimeout(() => setToastMessage(null), 3000);
      },
      onError: (err: any) => {
        setToastMessage({ type: "error", text: err.response?.data?.message || "Failed to save settings" });
        setTimeout(() => setToastMessage(null), 3000);
      },
    });
  };

  const handleDiscard = () => {
    if (settings) {
      setLocalSettings({ ...settings });
    }
    const rp = paymentSettingsList?.find((p) => p.provider === "razorpay");
    if (rp) {
      setRazorpayEnabled(rp.is_enabled);
    }
    const wt = paymentSettingsList?.find((p) => p.provider === "wallet");
    if (wt) {
      setWalletEnabled(wt.is_enabled);
    }
  };

  return (
    <AdminLayout title="Manage Platform">
      {toastMessage && (
        <div className={`fixed top-4 right-4 z-50 px-4 py-3 rounded-lg shadow-lg text-sm font-medium ${
          toastMessage.type === 'success' ? 'bg-green-100 text-green-800 border border-green-200' : 'bg-red-100 text-red-800 border border-red-200'
        }`}>
          {toastMessage.text}
        </div>
      )}
      
      <div className="flex flex-col gap-6 lg:flex-row items-start relative pb-24">
        
        {}
        <div className="flex-1 flex flex-col gap-6 w-full lg:max-w-3xl">

          {}
          <div className="bg-white rounded-[20px] p-6 shadow-sm border border-gray-100 h-full">
            <h3 className="text-lg font-bold text-[#111827] flex items-center gap-2 mb-6">
              <span className="w-8 h-8 rounded-full bg-blue-50 text-blue-500 flex items-center justify-center">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
              </span>
              User Management
            </h3>
            <div className="space-y-6">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-700">Enable user registrations</span>
                <Toggle isChecked={localSettings.enable_user_registration} onToggle={() => handleToggle("enable_user_registration")} />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-700 flex items-center gap-2">
                  <svg className="w-5 h-5 text-gray-500" viewBox="0 0 24 24" fill="currentcolor"><path d="M12.545,10.239v3.821h5.445c-0.712,2.315-2.647,3.972-5.445,3.972c-3.332,0-6.033-2.701-6.033-6.032s2.701-6.032,6.033-6.032c1.498,0,2.866,0.549,3.921,1.453l2.814-2.814C17.503,2.988,15.139,2,12.545,2C7.021,2,2.543,6.477,2.543,12s4.478,10,10.002,10c8.396,0,10.249-7.85,9.426-11.748L12.545,10.239z"/></svg>
                  Allow Google login
                </span>
                <Toggle isChecked={localSettings.allow_google_login} onToggle={() => handleToggle("allow_google_login")} />
              </div>
            </div>
          </div>

          {}
          <div className="bg-white rounded-[20px] p-6 shadow-sm border border-gray-100">
            <h3 className="text-lg font-bold text-[#111827] flex items-center gap-2 mb-6">
              <span className="w-8 h-8 rounded-full bg-teal-50 text-teal-500 flex items-center justify-center">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><polyline points="16 11 18 13 22 9"/></svg>
              </span>
              Organizer Control
            </h3>
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-gray-700">Require admin approval for new organizers</span>
              <Toggle isChecked={localSettings.require_admin_approval_for_organizers} onToggle={() => handleToggle("require_admin_approval_for_organizers")} />
            </div>
          </div>

          {}
          <div className="bg-white rounded-[20px] p-6 shadow-sm border border-gray-100 overflow-hidden">
             <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-bold text-[#111827] flex items-center gap-2">
                  <span className="w-8 h-8 rounded-full bg-[#1c2438] text-white flex items-center justify-center">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
                  </span>
                  Admin Management
                </h3>
                <button className="text-blue-600 text-sm font-semibold border border-blue-200 rounded-lg px-4 py-2 hover:bg-blue-50 transition">
                  + Add Admin
                </button>
             </div>
             
             <div className="overflow-x-auto">
               <table className="w-full text-sm text-left">
                 <thead className="bg-[#f8fafc] text-gray-700 text-xs uppercase font-bold rounded-t-xl overflow-hidden hidden sm:table-header-group">
                   <tr>
                     <th className="px-4 py-3 rounded-tl-xl whitespace-nowrap">Admin Name</th>
                     <th className="px-4 py-3 whitespace-nowrap">Admin Role</th>
                     <th className="px-4 py-3 whitespace-nowrap">Permissions</th>
                     <th className="px-4 py-3 whitespace-nowrap">Status</th>
                     <th className="px-4 py-3 rounded-tr-xl whitespace-nowrap text-right">Action</th>
                   </tr>
                 </thead>
                 <tbody className="divide-y divide-gray-100">
                   {adminsData?.admins?.map((admin: any) => (
                     <tr key={admin.id} className="hover:bg-gray-50 flex flex-col sm:table-row py-4 sm:py-0 border-b sm:border-b-0 border-gray-100">
                       <td className="px-4 py-3 font-medium text-gray-900 flex justify-between sm:table-cell">
                         <span className="sm:hidden font-bold">Name</span> {admin.name}
                       </td>
                       <td className="px-4 py-3 text-gray-600 flex justify-between sm:table-cell">
                         <span className="sm:hidden font-bold">Role</span> {admin.role}
                       </td>
                       <td className="px-4 py-3 text-gray-600 flex justify-between sm:table-cell">
                          <span className="sm:hidden font-bold">Perms</span> {admin.permissions}
                       </td>
                       <td className="px-4 py-3 flex justify-between items-center sm:table-cell">
                         <span className="sm:hidden font-bold">Status</span>
                         <span className={`px-3 py-1 rounded-full text-xs font-semibold ${admin.status === "Active" ? "bg-teal-100 text-teal-700" : "bg-rose-100 text-rose-700"}`}>
                           {admin.status}
                         </span>
                       </td>
                       <td className="px-4 py-3 text-right">
                         <button className="border border-gray-200 text-gray-700 px-3 py-1.5 rounded-lg text-xs font-semibold hover:bg-gray-100 transition inline-flex items-center gap-1 w-full sm:w-auto justify-center">
                           View
                           <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="6 9 12 15 18 9"/></svg>
                         </button>
                       </td>
                     </tr>
                   ))}
                   {!adminsData?.admins?.length && (
                      <tr>
                        <td colSpan={5} className="text-center py-6 text-gray-500">No admins found</td>
                      </tr>
                   )}
                 </tbody>
               </table>
             </div>
          </div>
        </div>

        {}
        <div className="flex flex-col gap-6 w-full lg:w-[400px]">

           {}
           <div className="bg-white rounded-[20px] p-6 shadow-sm border border-gray-100">
             <h3 className="text-lg font-bold text-[#111827] flex items-center gap-2 mb-6">
              <span className="w-8 h-8 rounded-full bg-slate-100 text-slate-600 flex items-center justify-center">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              </span>
              Event Policies
            </h3>
            
            <div className="space-y-6">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-700">Events require admin approval</span>
                <Toggle isChecked={localSettings.require_admin_approval_for_events} onToggle={() => handleToggle("require_admin_approval_for_events")} />
              </div>
              
              <div className="flex items-center justify-between">
                 <span className="text-sm font-medium text-gray-700 flex items-center gap-2">
                    <ShieldCheck className="w-4 h-4 text-blue-500" />
                    Refund Window
                 </span>
                 <div className="flex items-center gap-2">
                   <input 
                     type="number" 
                     className="w-12 h-8 border border-gray-200 rounded-lg text-center text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500"
                     value={localSettings.refund_window_days} 
                     onChange={(e) => handleChange("refund_window_days", parseInt(e.target.value) || 0)}
                     min={0}
                   />
                   <span className="text-xs text-gray-500 font-medium">days before</span>
                 </div>
              </div>

              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-700 flex items-center gap-2">
                  <span className="w-4 h-4 rounded-full bg-blue-500 text-white flex items-center justify-center text-[10px]">✓</span>
                  Allow event cancellation
                </span>
                <Toggle isChecked={localSettings.allow_event_cancellation} onToggle={() => handleToggle("allow_event_cancellation")} />
              </div>
            </div>
           </div>

           {}
           <div className="bg-white rounded-[20px] p-6 shadow-sm border border-gray-100">
             <h3 className="text-lg font-bold text-[#111827] flex items-center gap-2 mb-6">
              <span className="w-8 h-8 rounded-full bg-emerald-50 text-emerald-500 flex items-center justify-center">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="6" width="20" height="12" rx="2"/><circle cx="12" cy="12" r="2"/><path d="M6 12h.01M18 12h.01"/></svg>
              </span>
              Payment Settings
            </h3>

            <div className="space-y-6">
              <div className="flex items-center justify-between">
                 <span className="text-sm font-medium text-gray-700">Platform Fee (per ticket)</span>
                 
                 <div className="flex items-center gap-2">
                    <input 
                      type="number" 
                      className="w-16 h-8 border border-gray-200 rounded-lg text-center text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500"
                      value={localSettings.platform_fee_value === 0 ? "" : localSettings.platform_fee_value}
                      onChange={(e) => handleChange("platform_fee_value", parseFloat(e.target.value) || 0)}
                      step="0.01"
                      min="0"
                    />
                    <div className="flex bg-gray-100 rounded-md p-0.5">
                       <button 
                          className={`w-7 h-7 flex items-center justify-center rounded text-xs font-bold transition-colors ${localSettings.platform_fee_type === 'fixed' ? 'bg-white shadow text-gray-900' : 'text-gray-500'}`}
                          onClick={() => handleChange('platform_fee_type', 'fixed')}
                       >
                         ₹
                       </button>
                       <button 
                          className={`w-7 h-7 flex items-center justify-center rounded text-xs font-bold transition-colors ${localSettings.platform_fee_type === 'percentage' ? 'bg-white shadow text-gray-900' : 'text-gray-500'}`}
                          onClick={() => handleChange('platform_fee_type', 'percentage')}
                       >
                         %
                       </button>
                    </div>
                 </div>
              </div>

              <div>
                <span className="text-sm font-medium text-gray-700 block mb-3">Payment Modes & Gateways</span>
                <div className="border border-gray-100 rounded-xl bg-gray-50/50 p-2 space-y-1">
                   <div className="flex items-center justify-between p-2 pb-1">
                      <span className="text-sm text-gray-800 font-medium tracking-tight">Razorpay</span>
                      <Toggle isChecked={razorpayEnabled} onToggle={() => setRazorpayEnabled(!razorpayEnabled)} />
                   </div>
                   <div className="flex items-center justify-between p-2 pt-1 border-t border-gray-100">
                      <span className="text-sm text-gray-800 font-medium tracking-tight">Pay using Wallet</span>
                      <Toggle isChecked={walletEnabled} onToggle={() => setWalletEnabled(!walletEnabled)} />
                   </div>
                   {}
                </div>
              </div>
            </div>
           </div>

        </div>

      </div>
      
      {}
      <div className="fixed bottom-0 right-0 left-[240px] bg-white/80 backdrop-blur-md border-t border-gray-200 p-4 flex justify-end gap-3 z-20 px-8 shadow-[0_-10px_20px_-10px_rgba(0,0,0,0.05)]">
        <button 
          onClick={handleDiscard}
          className="px-6 py-2.5 rounded-xl border border-gray-200 text-gray-600 font-bold text-sm hover:bg-gray-50 transition"
        >
          Discard Changes
        </button>
        <button 
          onClick={handleSave}
          disabled={isSaving}
          className="px-6 py-2.5 rounded-xl bg-[#5d779f] text-white font-bold text-sm hover:bg-[#4a6184] transition shadow-[0_4px_12px_rgba(93,119,159,0.3)] disabled:opacity-70 flex items-center gap-2"
        >
          {isSaving ? <Loader2 className="w-4 h-4 animate-spin" /> : "Save Changes"}
        </button>
      </div>

    </AdminLayout>
  );
}

function Toggle({ isChecked, onToggle, disabled = false }: { isChecked: boolean, onToggle: () => void, disabled?: boolean }) {
  return (
    <button 
      type="button"
      role="switch" 
      aria-checked={isChecked}
      disabled={disabled}
      onClick={onToggle}
      className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-opacity-75 ${
        isChecked ? 'bg-[#5d779f]' : 'bg-gray-200'
      } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
    >
      <span className="sr-only">Toggle switch</span>
      <span
        aria-hidden="true"
        className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
          isChecked ? 'translate-x-5' : 'translate-x-0'
        }`}
      />
    </button>
  );
}
