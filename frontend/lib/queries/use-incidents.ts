import { useQuery } from "@tanstack/react-query";

import {
  listIncidents,
  listMonitorIncidents,
} from "@/lib/api/incidents";
import { queryKeys } from "@/lib/queries/query-keys";

export function useIncidents() {
  return useQuery({
    queryFn: () => listIncidents(),
    queryKey: queryKeys.incidents,
  });
}

export function useMonitorIncidents(monitorId: string) {
  return useQuery({
    enabled: Boolean(monitorId),
    queryFn: () => listMonitorIncidents(monitorId),
    queryKey: queryKeys.monitorIncidents(monitorId),
  });
}
