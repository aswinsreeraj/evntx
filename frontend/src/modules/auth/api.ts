import api from "../../services/axios";
import { tokenManager } from "../../services/tokenManager";
import { useAuthStore } from "./store/authStore";

type VerifyOtpResponse = {
    access_token: string;
    refresh_token: string;
    expires_in: number;
    user: {
        id: string;
        name: string;
        roles: string[];
    };
};

export const authApi = {
    async requestOtp(email: string) {
        console.log("Requesting OTP for email:", email);
        return api.post("/auth/otp/request", { email });
    },

    async verifyOtp(email: string, otp: string) {
        const response = await api.post("/auth/otp/verify", {
            email,
            otp,
        });

        const data: VerifyOtpResponse = response.data.data;

        // Set access token in memory
        tokenManager.setToken(data.access_token);

        // Update auth store
        useAuthStore.getState().setAuth(
            {
                id: data.user.id,
                name: data.user.name,
            },
            data.user.roles
        );

        return response;
    },

    async googleLogin(idToken: string) {
        const response = await api.post("/auth/oauth/google", {
            id_token: idToken,
        });

        const data = response.data.data;

        tokenManager.setToken(data.access_token);

        // Need user details depending on backend response
        // If backend doesn't return user object here, fetch /users/me later

        return response;
    },

    async logout() {
        await api.post("/auth/logout");

        tokenManager.clearToken();
        useAuthStore.getState().logout();
    },
};