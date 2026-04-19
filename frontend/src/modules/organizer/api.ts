import api from "../../services/axios";
import type { PayoutCredentialPayload, PayoutsResponse } from "../user/api";

export interface PlatformSettings {
  allow_event_cancellation: boolean;
  enable_user_registration: boolean;
  allow_google_login: boolean;
  require_admin_approval_for_organizers: boolean;
  require_admin_approval_for_events: boolean;
  refund_window_days: number;
  platform_fee_value: number;
  platform_fee_type: "fixed" | "percentage";
}

export const organizerWalletSummaryQueryKey = ["organizer-wallet-summary"] as const;

export interface TicketInput {
  id?: string;
  name: string;
  price: number;
  total_quantity: number;
}

export interface PersonnelInput {
  id?: string;
  name: string;
  role: string;
  image?: string;
  profile_link?: string;
}

export interface EventDetailsInput {
  description: string;
  venue_address: string;
  map_url?: string;
  total_capacity: number;
  terms_and_conditions?: string;
}

export interface CreateEventPayload {
  title: string;
  city: string;
  venue_name: string;
  category?: string;
  start_time: string;
  end_time: string;
  tags?: string[];
  cover_image_url?: string;
  status?: string;
  details: EventDetailsInput;
  ticket_types: TicketInput[];
  key_personnel?: PersonnelInput[];
}

export interface UpdateEventPayload {
  title?: string;
  city?: string;
  venue_name?: string;
  category?: string;
  start_time?: string;
  end_time?: string;
  tags?: string[];
  cover_image_url?: string;
  status?: string;
  details?: Partial<EventDetailsInput>;
  ticket_types?: TicketInput[];
  key_personnel?: PersonnelInput[];
}

export interface OrganizerWalletSummary {
  available_balance: number;
  reserve_balance: number;
  pending_balance: number;
  total_credited: number;
  total_debited: number;
}


export interface StatCardData {
  value: number;
  percentage: number;
  label?: string;
}

export interface RevenuePointData {
  date: string;
  amount: number;
}

export interface EventSalesBreakdownData {
  name: string;
  value: number;
}

export interface OrganizerDashboardStats {
  total_revenue: StatCardData;
  tickets_sold: StatCardData;
  active_events: StatCardData;
  pending_events: StatCardData;
  revenue_overview: RevenuePointData[];
  sales_breakdown: EventSalesBreakdownData[];
}

export interface TicketSalesProportionData {
  name: string;
  tickets_sold: number;
  percentage_total: number;
}

export interface SalesReportStats {
  total_revenue: StatCardData;
  tickets_sold: StatCardData;
  revenue_over_time: RevenuePointData[];
  tickets_per_event: TicketSalesProportionData[];
}


export interface FunnelStepData {
  label: string;
  count: number;
  percentage: number;
}

export interface PeakUsagePointData {
  label: string;
  viewing: number;
  checkout: number;
}

export interface EngagementReportStats {
  page_views: StatCardData;
  conversion_rate: StatCardData;
  user_journey: FunnelStepData[];
  peak_usage: PeakUsagePointData[];
}

export interface CheckInResponse {
  ticket_id: string;
  ticket_code: string;
  status: string;
  checked_in_at: string;
}

export const organizerApi = {
  async createEvent(payload: CreateEventPayload) {
    const res = await api.post("/organizer/events", payload);
    return res.data;
  },

  async updateEvent(eventId: string, payload: UpdateEventPayload) {
    const res = await api.put(`/organizer/events/${eventId}`, payload);
    return res.data;
  },

  async getDashboardStats(): Promise<OrganizerDashboardStats> {
    const res = await api.get(`/organizer/dashboard`);
    return res.data.data;
  },

  async getSalesReportStats(eventId?: string, startDate?: string, endDate?: string): Promise<SalesReportStats> {
    const params = new URLSearchParams();
    if (eventId) params.append("event_id", eventId);
    if (startDate) params.append("start_date", startDate);
    if (endDate) params.append("end_date", endDate);
    const res = await api.get(`/organizer/reports/sales?${params.toString()}`);
    return res.data.data;
  },

  async getEngagementReportStats(eventId?: string, startDate?: string, endDate?: string): Promise<EngagementReportStats> {
    const params = new URLSearchParams();
    if (eventId) params.append("event_id", eventId);
    if (startDate) params.append("start_date", startDate);
    if (endDate) params.append("end_date", endDate);
    const res = await api.get(`/organizer/reports/engagement?${params.toString()}`);
    return res.data.data;
  },

  async getOrganizerEvents(status?: string) {
    const res = await api.get(`/organizer/events`, { params: { status } });
    return res.data.data;
  },

  async getWalletSummary(): Promise<OrganizerWalletSummary> {
    const res = await api.get("/organizer/wallet");
    return res.data.data;
  },

  async addPayoutCredentials(payload: PayoutCredentialPayload) {
    const res = await api.post("/organizer/payout/credentials", payload);
    return res.data;
  },

  async requestPayout(amount: number) {
    const res = await api.post("/organizer/wallet/payout", { amount });
    return res.data;
  },

  async getPayouts(): Promise<PayoutsResponse> {
    const res = await api.get("/organizer/payouts");
    return res.data.data;
  },

  async getEventBySlug(slug: string) {
    const res = await api.get(`/organizer/events/slug/${slug}`);
    return res.data.data;
  },

  async deleteEvent(eventId: string) {
    const res = await api.delete(`/organizer/events/${eventId}`);
    return res.data;
  },

  async requestEventCancellation(eventId: string, reason: string) {
    const res = await api.post(`/organizer/events/${eventId}/cancel-request`, { reason });
    return res.data;
  },

  async submitEvent(eventId: string) {
    const res = await api.post(`/organizer/events/${eventId}/submit`);
    return res.data;
  },

  async checkInTicket(eventId: string, ticketCode: string): Promise<CheckInResponse> {
    const res = await api.post(`/events/${eventId}/check-in`, {
      ticket_code: ticketCode,
    });
    return res.data.data;
  },

  async getPlatformSettings(): Promise<PlatformSettings> {
    const res = await api.get("/settings");
    return res.data.data;
  },

  async uploadImage(file: File) {
    const formData = new FormData();
    formData.append("image", file);
    const res = await api.post("/organizer/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return res.data.data;
  },
};
