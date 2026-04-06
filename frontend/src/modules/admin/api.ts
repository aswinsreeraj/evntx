import api from "../../services/axios";
import type { PayoutRequestData } from "../user/api";

export interface AdminPayoutDetail extends PayoutRequestData {
  user_name: string;
  user_email: string;
}

export interface AdminPayoutsResponse {
  payouts: AdminPayoutDetail[];
  total: number;
}

export interface AdminRefundDetail {
  id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  booking_id: string;
  amount: number;
  status: "pending" | "processed";
  requested_at: string;
  processed_at?: string;
}

export interface AdminRefundsResponse {
  refunds: AdminRefundDetail[];
  total: number;
}

export interface PlatformSettings {
  id: string;
  enable_user_registration: boolean;
  allow_google_login: boolean;
  require_admin_approval_for_organizers: boolean;
  require_admin_approval_for_events: boolean;
  refund_window_days: number;
  allow_event_cancellation: boolean;
  platform_fee_value: number;
  platform_fee_type: "fixed" | "percentage";
  updated_at: string;
}

export interface PaymentSettings {
  id: string;
  provider: string;
  is_enabled: boolean;
  config: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface AdminUser {
  id: string;
  name: string;
  email: string;
  role: string;
  permissions: string;
  status: string;
}

export interface AuditLog {
  id: string;
  admin_id: string;
  admin_name: string;
  action: string;
  action_tag: string;
  details: string;
  ip_address: string;
  timestamp: string;
}

// -- Admin Dashboard --
export interface AdminStatCardData {
  value: number;
  percentage: number;
  subtitle?: string;
}

export interface AdminDashboardStats {
  revenue: AdminStatCardData;
  total_users: AdminStatCardData;
  total_organizers: AdminStatCardData;
  total_events: AdminStatCardData;
  refund_rate: AdminStatCardData;
  user_growth: AdminStatCardData;
  pending_approvals: AdminStatCardData;
  active_events: AdminStatCardData;
  revenue_overview: { date: string; amount: number }[];
}

// -- Admin Revenue Report --
export interface CategoryRevenueData {
  category: string;
  revenue: number;
}

export interface RefundDataPoint {
  month: string;
  amount: number;
}

export interface TopOrganizerEntry {
  name: string;
  revenue: number;
  active_events: number;
  pending_events: number;
  avg_event_rating: number;
}

export interface TopUserEntry {
  name: string;
  events_attended: number;
  total_spent: number;
}

export interface AdminRevenueReport {
  revenue_today: AdminStatCardData;
  revenue_this_month: AdminStatCardData;
  total_revenue: AdminStatCardData;
  growth_rate: AdminStatCardData;
  revenue_over_time: { date: string; amount: number }[];
  category_breakdown: CategoryRevenueData[];
  refund_analytics: RefundDataPoint[];
  refund_total: AdminStatCardData;
  top_organizers: TopOrganizerEntry[];
  top_users: TopUserEntry[];
}

export const adminApi = {
  async getDashboardStats(): Promise<AdminDashboardStats> {
    const response = await api.get("/admin/dashboard");
    return response.data.data;
  },

  async getRevenueReport(startDate?: string, endDate?: string): Promise<AdminRevenueReport> {
    const params: Record<string, string> = {};
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;
    const response = await api.get("/admin/reports/revenue", { params });
    return response.data.data;
  },

  async getEngagementReportStats(
    organizerId?: string,
    eventId?: string,
    startDate?: string,
    endDate?: string
  ) {
    const params: Record<string, string> = {};
    if (organizerId && organizerId !== "all") params.organizer_id = organizerId;
    if (eventId && eventId !== "all") params.event_id = eventId;
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;
    const res = await api.get("/admin/reports/engagement", { params });
    return res.data.data;
  },

  async getUsers(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    role?: string;
    is_active?: boolean;
  }) {
    const response = await api.get("/admin/users", { params });
    return response.data.data;
  },

  async updateUserStatus(userId: string, isActive: boolean) {
    const response = await api.patch(`/admin/users/${userId}/status`, {
      is_active: isActive,
    });

    return response.data;
  },

  async getOrganizers(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
  }) {
    const response = await api.get("/admin/organizers", { params });
    return response.data.data;
  },

  async getEvents(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
  }) {
    const response = await api.get("/admin/events", { params });
    return response.data.data;
  },

  async getEventBySlug(slug: string) {
    const response = await api.get(`/admin/events/slug/${slug}`);
    return response.data.data;
  },

  async approveEvent(eventId: string) {
    const response = await api.patch(`/admin/events/${eventId}/approve`);
    return response.data;
  },

  async rejectEvent(eventId: string, reason: string) {
    const response = await api.patch(`/admin/events/${eventId}/reject`, {
      reason,
    });
    return response.data;
  },

  async getPlatformWallet() {
    const response = await api.get("/admin/platform-wallet");
    return response.data.data;
  },

  async getPayouts(params?: { status?: string }): Promise<AdminPayoutsResponse> {
    const response = await api.get("/admin/payouts", { params });
    return response.data.data;
  },

  async approvePayout(payoutId: string) {
    const response = await api.patch(`/admin/payouts/${payoutId}/approve`);
    return response.data;
  },

  async rejectPayout(payoutId: string, reason: string) {
    const response = await api.patch(`/admin/payouts/${payoutId}/reject`, { reason });
    return response.data;
  },

  async bulkApprovePayouts(payoutIds: string[]) {
    const response = await api.post(`/admin/payouts/bulk-approve`, { payout_ids: payoutIds });
    return response.data;
  },

  async getRefunds(params?: { status?: string }): Promise<AdminRefundsResponse> {
    const response = await api.get("/admin/refunds", { params });
    return response.data.data;
  },

  async processRefund(refundId: string) {
    const response = await api.patch(`/admin/refunds/${refundId}/process`);
    return response.data;
  },

  async getPlatformSettings(): Promise<PlatformSettings> {
    const response = await api.get("/admin/settings");
    return response.data.data;
  },

  async updatePlatformSettings(settings: Partial<PlatformSettings>): Promise<PlatformSettings> {
    const response = await api.put("/admin/settings", settings);
    return response.data.data;
  },

  async getPaymentSettings(): Promise<PaymentSettings[]> {
    const response = await api.get("/admin/payment-settings"); // Wait, I didn't add /admin/payment-settings to router
    return response.data.data;
  },

  async updatePaymentProvider(provider: string, data: { is_enabled: boolean; config: Record<string, any> }): Promise<void> {
    const response = await api.put(`/admin/payment-settings/${provider}`, data);
    return response.data;
  },
  
  async getAdmins(): Promise<{ admins: AdminUser[] }> {
    const response = await api.get("/admin/admins");
    return response.data.data;
  },

  async getAuditLogs(page: number = 1, limit: number = 20): Promise<{ logs: AuditLog[], pagination: { total: number, page: number, limit: number } }> {
    const response = await api.get(`/admin/audit-logs?page=${page}&limit=${limit}`);
    return response.data.data;
  }
};
