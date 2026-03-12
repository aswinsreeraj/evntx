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

    return response.data.data;
  },
};