import Layout from "../shared/Layout";
import HomePage from "../modules/auth/pages/HomePage";
import LoginPage from "../modules/auth/pages/LoginPage";
import VerifyOtpPage from "../modules/auth/pages/VerifyOtpPage";
import ProfilePage from "../modules/user/pages/ProfilePage";
import ProtectedRoute from "../shared/components/ProtectedRoute";
import { createBrowserRouter } from "react-router-dom";

export const router = createBrowserRouter([
    {
        element: <Layout />,
        children: [
            { path: "/", element: <HomePage /> },
            { path: "/login", element: <LoginPage /> },
            { path: "/verify-otp", element: <VerifyOtpPage /> },
            {   path: "/profile",
                element: (
                    <ProtectedRoute>
                        <ProfilePage />
                    </ProtectedRoute>
                ),  
            },
        ],
    },
]);