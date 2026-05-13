import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  createMonitor,
  getMonitor,
  listMonitors,
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
