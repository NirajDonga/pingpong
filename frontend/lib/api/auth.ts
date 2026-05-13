import { api, setToken } from "@/lib/api/client";
import type { AuthResponse } from "@/lib/api/types";

type AuthInput = {
  email: string;
  password: string;
};

export async function login(input: AuthInput) {
  const response = await api<AuthResponse>("/login", {
    body: input,
    method: "POST",
  });
  setToken(response.token);
  return response;
}

export async function register(input: AuthInput) {
  const response = await api<AuthResponse>("/register", {
    body: input,
    method: "POST",
  });
  setToken(response.token);
  return response;
}
