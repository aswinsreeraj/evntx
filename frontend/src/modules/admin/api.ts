import api from "../../services/axios";

export const adminApi = {
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
};
