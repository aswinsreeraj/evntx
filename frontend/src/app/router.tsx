import Layout from "../shared/components/Layout";
import HomePage from "../modules/home/pages/HomePage";
import ProfilePage from "../modules/user/pages/ProfilePage";
import ProtectedRoute from "../shared/components/ProtectedRoute";
import { createBrowserRouter } from "react-router-dom";
import EventListPage from "../modules/events/pages/EventListPage";
import EventDetailPage from "../modules/events/pages/EventDetailPage";
import AdminLoginPage from "../modules/admin/pages/AdminLoginPage";
import UserManagementPage from "../modules/admin/pages/UserManagementPage";
import OrganizerProfile from "../modules/organizer/pages/Profile";

export const router = createBrowserRouter([
    {
        element: <Layout />,
        children: [
            { path: "/", element: <HomePage /> },
            { path: "/login", element: <HomePage /> },
            {   path: "/profile",
                element: (
                    <ProtectedRoute>
                        <ProfilePage />
                    </ProtectedRoute>
                ),  
            },
            { path: "/events", element: <EventListPage /> },
            { path: "/events/:eventId", element: <EventDetailPage /> },
        ],
    },
    {
        path: "/organizer/profile",
        element: (
            <ProtectedRoute roles={["organizer"]}>
                <OrganizerProfile />
            </ProtectedRoute>
        ),
    },
    { path: "/admin/login", element: <AdminLoginPage /> },
    {
        path: "/admin/users",
        element: (
            <ProtectedRoute roles={["admin"]}>
                <UserManagementPage />
            </ProtectedRoute>
        ),
    }
]);