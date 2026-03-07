import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { router } from "./app/router";
import { AppProviders } from "./app/providers";
import { GoogleOAuthProvider } from "@react-oauth/google";
import "./app/styles.css"

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <GoogleOAuthProvider clientId={import.meta.env.VITE_GOOGLE_CLIENT_ID || ""}>
      <AppProviders>
        <RouterProvider router={router} />
      </AppProviders>
    </GoogleOAuthProvider>
  </StrictMode>
);