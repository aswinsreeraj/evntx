import { createBrowserRouter } from "react-router-dom";
import HomePage from "../modules/auth/pages/HomePage";
import LoginPage from "../modules/auth/pages/LoginPage";
import Layout from "../shared/Layout";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <Layout />,
        children: [
            {
                index: true,
                element: <HomePage />,
            },
            {
                path: "login",
                element: <LoginPage />,
            },
        ],
    },
]);