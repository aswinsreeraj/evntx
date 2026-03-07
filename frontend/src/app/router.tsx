import Layout from "../shared/components/Layout";
import HomePage from "../modules/home/pages/HomePage";
import LoginPage from "../modules/auth/pages/LoginPage";
import VerifyOtpPage from "../modules/auth/pages/VerifyOtpPage";
import ProfilePage from "../modules/user/pages/ProfilePage";
import ProtectedRoute from "../shared/components/ProtectedRoute";
import { createBrowserRouter } from "react-router-dom";
import EventListPage from "../modules/events/pages/EventListPage";
import EventDetailPage from "../modules/events/pages/EventDetailPage";
import AdminLoginPage from "../modules/admin/pages/AdminLoginPage";
import UserManagementPage from "../modules/admin/pages/UserManagementPage";

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
            { path: "/events", element: <EventListPage /> },
            { path: "/events/:eventId", element: <EventDetailPage /> },
            { path: "/admin/login", element: <AdminLoginPage /> },
            {
            path: "/admin/users",
            element: (
                <ProtectedRoute roles={["admin"]}>
                <UserManagementPage />
                </ProtectedRoute>
            ),
            }
        ],
    },
]);