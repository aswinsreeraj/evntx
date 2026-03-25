import axios from "axios";
import { tokenManager } from "./tokenManager";
import { useAuthStore } from "../modules/auth/store/authStore";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    withCredentials: true,
});

api.interceptors.request.use((config) => {
    const token = tokenManager.getToken();

    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
});

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                const refreshToken = tokenManager.getRefreshToken();
                if (!refreshToken) {
                    throw new Error("No refresh token available");
                }

                const response = await axios.post(
                    `${import.meta.env.VITE_API_BASE_URL}/auth/refresh`,
                    { refresh_token: refreshToken },
                    { withCredentials: true }
                );

                const newToken = response.data.data.access_token;
                tokenManager.setToken(newToken);

                return api(originalRequest);
            } catch (refreshError) {
                tokenManager.clearToken();
                useAuthStore.getState().logout();
                return Promise.reject(refreshError);
            }
        }

        return Promise.reject(error);
    }
);

export default api;