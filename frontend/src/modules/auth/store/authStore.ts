import { create } from "zustand";
import { persist } from "zustand/middleware";

type User = {
    id: string;
    name: string;
};

type AuthState = {
    user: User | null;
    roles: string[];
    isAuthenticated: boolean;
    authModalOpen: "goer" | "organizer" | null;
    setAuth: (user: User, roles: string[]) => void;
    logout: () => void;
    openAuthModal: (mode: "goer" | "organizer") => void;
    closeAuthModal: () => void;
};

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            roles: [],
            isAuthenticated: false,
            authModalOpen: null,

            openAuthModal: (mode) => set({ authModalOpen: mode }),
            closeAuthModal: () => set({ authModalOpen: null }),

            setAuth: (user, roles) =>
                set({
                    user,
                    roles,
                    isAuthenticated: true,
                }),

            logout: () =>
                set({
                    user: null,
                    roles: [],
                    isAuthenticated: false,
                }),
        }),
        {
            name: "auth-storage",
        }
    )
);