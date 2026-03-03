import Layout from "../shared/components/Layout";
import HomePage from "../modules/auth/pages/HomePage";
import LoginPage from "../modules/auth/pages/LoginPage";
import VerifyOtpPage from "../modules/auth/pages/VerifyOtpPage";

export const router = createBrowserRouter([
    {
        element: <Layout />,
        children: [
            { path: "/", element: <HomePage /> },
            { path: "/login", element: <LoginPage /> },
            { path: "/verify-otp", element: <VerifyOtpPage /> },
        ],
    },
]);