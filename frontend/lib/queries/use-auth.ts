import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { currentUser, login, logout, register } from "@/lib/api/auth";
import { queryKeys } from "@/lib/queries/query-keys";

export function useCurrentUser() {
  return useQuery({
    queryFn: currentUser,
    queryKey: queryKeys.currentUser,
    retry: false,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: login,
    onSuccess: (user) => {
      queryClient.setQueryData(queryKeys.currentUser, { id: user.id });
    },
  });
}

export function useRegister() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: register,
    onSuccess: (user) => {
      queryClient.setQueryData(queryKeys.currentUser, { id: user.id });
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.setQueryData(queryKeys.currentUser, null);
    },
  });
}
