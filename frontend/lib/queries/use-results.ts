import { useQuery } from "@tanstack/react-query";

import { listMonitorChecks } from "@/lib/api/results";
import { queryKeys } from "@/lib/queries/query-keys";

export function useMonitorChecks(monitorId: string) {
  return useQuery({
    enabled: Boolean(monitorId),
    queryFn: () => listMonitorChecks(monitorId),
    queryKey: queryKeys.monitorChecks(monitorId),
    refetchInterval: 30_000,
  });
}
