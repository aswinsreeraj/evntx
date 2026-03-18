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
        const response = await api.post("/auth/otp/request", { email });
        return response.data; 
    },

    async verifyOtp(email: string, otp: string, name?: string) {
        const response = await api.post("/auth/otp/verify", {
            email,
            otp,
            name,
        });

        const data: VerifyOtpResponse = response.data.data;

        
        tokenManager.setToken(data.access_token);
        tokenManager.setRefreshToken(data.refresh_token);

        
        useAuthStore.getState().setAuth(
            {
                id: data.user.id,
                name: data.user.name,
            },
            data.user.roles
        );

        return response;
    },

    async register(email: string, otp: string, name: string, dob: string, gender: string, role?: string, organization_name?: string) {
        const response = await api.post("/auth/register", {
            email,
            otp,
            name,
            dob,
            gender,
            role,
            organization_name,
        });

        const data: VerifyOtpResponse = response.data.data;

        tokenManager.setToken(data.access_token);
        tokenManager.setRefreshToken(data.refresh_token);

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
        tokenManager.setRefreshToken(data.refresh_token);

        try {
            const meResponse = await api.get("/users/me");
            const user = meResponse.data.data;
            useAuthStore.getState().setAuth(
                {
                    id: user.id,
                    name: user.name,
                },
                user.roles || []
            );
        } catch (err) {
            console.error("Failed to fetch user profile after Google login", err);
        }

        return response;
    },

    async logout() {
        await api.post("/auth/logout");

        tokenManager.clearToken();
        useAuthStore.getState().logout();
    },
};