import Layout from "../shared/components/Layout";
import HomePage from "../modules/home/pages/HomePage";
import ProfilePage from "../modules/user/pages/ProfilePage";
import MyBookingsPage from "../modules/user/pages/MyBookingsPage";
import CalendarPage from "../modules/user/pages/CalendarPage";
import ProtectedRoute from "../shared/components/ProtectedRoute";
import { createBrowserRouter } from "react-router-dom";
import EventListPage from "../modules/events/pages/EventListPage";
import EventDetailPage from "../modules/events/pages/EventDetailPage";
import EventBookingPage from "../modules/events/pages/EventBookingPage";
import BookingConfirmationPage from "../modules/events/pages/BookingConfirmationPage";
import AdminLoginPage from "../modules/admin/pages/AdminLoginPage";
import UserManagementPage from "../modules/admin/pages/UserManagementPage";
import OrganizerManagementPage from "../modules/admin/pages/OrganizerManagementPage";
import EventManagementPage from "../modules/admin/pages/EventManagementPage";
import OrganizerProfile from "../modules/organizer/pages/Profile";
import EventForm from "../modules/organizer/pages/EventForm";
import MyEvents from "../modules/organizer/pages/MyEvents";

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
            {
                path: "/events/:eventId/confirmation",
                element: (
                    <ProtectedRoute>
                        <BookingConfirmationPage />
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
    }
]);
