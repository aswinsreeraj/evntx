import api from "../../services/axios";

export type UpdateProfilePayload = {
  name: string;
  mobile: string;
  dob: string;
  gender: string;
  locations: string[];
};

export const userApi = {
  async getProfile() {
    const res = await api.get("/users/me");
    return res.data.data;
  },

  async updateProfile(payload: UpdateProfilePayload) {
    const res = await api.put("/users/me", payload);
    return res.data.data;
  },
};