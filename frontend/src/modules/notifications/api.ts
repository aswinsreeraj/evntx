import api from "../../services/axios";

export type NotificationItem = {
  id: string;
  user_id: string;
  type: string;
  title: string;
  message: string;
  is_read: boolean;
  metadata?: Record<string, unknown> | null;
  created_at: string;
};

export type NotificationsResponse = {
  notifications: NotificationItem[];
  unread_count: number;
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
};

export const notificationsApi = {
  async getNotifications(page = 1, limit = 6): Promise<NotificationsResponse> {
    const res = await api.get("/notifications", {
      params: { page, limit },
    });
    return res.data.data;
  },

  async markAsRead(notificationId: string): Promise<void> {
    await api.patch(`/notifications/${notificationId}/read`);
  },

  async markAllAsRead(): Promise<void> {
    await api.patch("/notifications/read-all");
  },

  async clearAll(): Promise<void> {
    await api.delete("/notifications");
  },
};
