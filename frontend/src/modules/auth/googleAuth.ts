import { authApi } from "./api";

export async function handleGoogleLogin(idToken: string) {
  const response = await authApi.googleLogin(idToken);
  return response;
}
