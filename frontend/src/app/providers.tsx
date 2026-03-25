import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import AuthHydrator from "./AuthHydrator";

const queryClient = new QueryClient();

export function AppProviders({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthHydrator />
      {children}
    </QueryClientProvider>
  );
}
