import { createBrowserRouter } from "react-router-dom";
import HomePage from "../modules/auth/pages/HomePage";
import LoginPage from "../modules/auth/pages/LoginPage";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <HomePage />,
    },
    {
        path: "/login",
        element: <LoginPage />,
    },
]);