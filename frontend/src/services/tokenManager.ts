const ACCESS_TOKEN_KEY = "auth_access_token";
const REFRESH_TOKEN_KEY = "auth_refresh_token";

export const tokenManager = {
    setToken(token: string) {
        localStorage.setItem(ACCESS_TOKEN_KEY, token);
    },

    getToken() {
        return localStorage.getItem(ACCESS_TOKEN_KEY);
    },

    setRefreshToken(token: string) {
        localStorage.setItem(REFRESH_TOKEN_KEY, token);
    },

    getRefreshToken() {
        return localStorage.getItem(REFRESH_TOKEN_KEY);
    },

    clearToken() {
        localStorage.removeItem(ACCESS_TOKEN_KEY);
        localStorage.removeItem(REFRESH_TOKEN_KEY);
    },
};