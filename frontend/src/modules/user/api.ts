import api from "../../services/axios";

export const userApi = {
  async getProfile() {
    const res = await api.get("/users/me");
    return res.data.data;
  },

  async updateProfile(payload: { name: string }) {
    const res = await api.put("/users/me", payload);
    return res.data.data;
  },
};