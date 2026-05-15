import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  createMonitor,
  deleteMonitor,
  getMonitor,
  listMonitors,
  setMonitorEnabled,
} from "@/lib/api/monitors";
import { queryKeys } from "@/lib/queries/query-keys";

export function useMonitors() {
  return useQuery({
    queryFn: listMonitors,
    queryKey: queryKeys.monitors,
  });
}

export function useMonitor(id: string) {
  return useQuery({
    enabled: Boolean(id),
    queryFn: () => getMonitor(id),
    queryKey: queryKeys.monitor(id),
  });
}

export function useCreateMonitor() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createMonitor,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.monitors });
    },
  });
}

export function useSetMonitorEnabled(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (enabled: boolean) => setMonitorEnabled(id, enabled),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.monitor(id) });
      queryClient.invalidateQueries({ queryKey: queryKeys.monitors });
    },
  });
}

export function useDeleteMonitor(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => deleteMonitor(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.monitors });
    },
  });
}
