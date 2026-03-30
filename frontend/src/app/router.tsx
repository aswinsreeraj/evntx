import Layout from "../shared/components/Layout";
import HomePage from "../modules/home/pages/HomePage";
import ProfilePage from "../modules/user/pages/ProfilePage";
import MyBookingsPage from "../modules/user/pages/MyBookingsPage";
import CalendarPage from "../modules/user/pages/CalendarPage";
import WalletPage from "../modules/user/pages/WalletPage";
import ProtectedRoute from "../shared/components/ProtectedRoute";
import { createBrowserRouter } from "react-router-dom";
import EventListPage from "../modules/events/pages/EventListPage";
import EventDetailPage from "../modules/events/pages/EventDetailPage";
import EventBookingPage from "../modules/events/pages/EventBookingPage";
import AdminLoginPage from "../modules/admin/pages/AdminLoginPage";
import UserManagementPage from "../modules/admin/pages/UserManagementPage";
import OrganizerManagementPage from "../modules/admin/pages/OrganizerManagementPage";
import EventManagementPage from "../modules/admin/pages/EventManagementPage";
import OrganizerProfile from "../modules/organizer/pages/Profile";
import EventForm from "../modules/organizer/pages/EventForm";
import MyEvents from "../modules/organizer/pages/MyEvents";
import OrganizerWalletPage from "../modules/organizer/pages/WalletPage";
import OrganizerCheckInPage from "../modules/organizer/pages/CheckInPage";
import PlatformWalletPage from "../modules/admin/pages/PlatformWalletPage";

export const router = createBrowserRouter([
    {
        element: <Layout />,
        children: [
            { path: "/", element: <HomePage /> },
            { path: "/login", element: <HomePage /> },
            {   path: "/profile",
                element: (
                    <ProtectedRoute roles={["goer"]}>
                        <ProfilePage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/profile/bookings",
                element: (
                    <ProtectedRoute roles={["goer"]}>
                        <MyBookingsPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/profile/calendar",
                element: (
                    <ProtectedRoute roles={["goer"]}>
                        <CalendarPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/wallet",
                element: (
                    <ProtectedRoute>
                        <WalletPage />
                    </ProtectedRoute>
                ),
            },
            { path: "/events", element: <EventListPage /> },
            { path: "/events/:eventId", element: <EventDetailPage /> },
            {
                path: "/events/:eventId/book",
                element: (
                    <ProtectedRoute>
                        <EventBookingPage />
                    </ProtectedRoute>
                ),
            },
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
    {
        path: "/organizer/events/create",
        element: (
            <ProtectedRoute roles={["organizer"]}>
                <EventForm />
            </ProtectedRoute>
        ),
    },
    {
        path: "/organizer/events/:eventId/edit",
        element: (
            <ProtectedRoute roles={["organizer"]}>
                <EventForm />
            </ProtectedRoute>
        ),
    },
    {
        path: "/organizer/events",
        element: (
            <ProtectedRoute roles={["organizer"]}>
                <MyEvents />
            </ProtectedRoute>
        ),
    },
    {
        path: "/organizer/events/:eventId/check-in",
        element: (
            <ProtectedRoute roles={["organizer"]}>
                <OrganizerCheckInPage />
            </ProtectedRoute>
        ),
    },
    {
        path: "/organizer/wallet",
        element: (
            <ProtectedRoute roles={["organizer"]}>
                <OrganizerWalletPage />
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
    },
    {
        path: "/admin/organizers",
        element: (
            <ProtectedRoute roles={["admin"]}>
                <OrganizerManagementPage />
            </ProtectedRoute>
        ),
    },
    {
        path: "/admin/events",
        element: (
            <ProtectedRoute roles={["admin"]}>
                <EventManagementPage />
            </ProtectedRoute>
        ),
    },
    {
        path: "/admin/platform-wallet",
        element: (
            <ProtectedRoute roles={["admin"]}>
                <PlatformWalletPage />
            </ProtectedRoute>
        ),
    }
]);
