import { ApiError, api } from "@/lib/api/client";
import type { User } from "@/lib/api/types";

type AuthInput = {
  email: string;
  password: string;
};

export type Session = {
  id: string;
};

export async function login(input: AuthInput) {
  return api<User>("/login", {
    body: input,
    method: "POST",
  });
}

export async function register(input: AuthInput) {
  return api<User>("/register", {
    body: input,
    method: "POST",
  });
}

export async function logout() {
  return api<void>("/logout", {
    method: "POST",
  });
}

export async function currentUser() {
  try {
    return await api<Session>("/me");
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null;
    }

    throw error;
  }
}
