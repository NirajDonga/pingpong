import { useEffect } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";

import { listMonitorChecks } from "@/lib/api/results";
import type { CheckResult } from "@/lib/api/types";
import { queryKeys } from "@/lib/queries/query-keys";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api";

function getWsUrl(monitorId: string): string {
  // Convert http(s):// to ws(s)://
  const wsBase = API_BASE_URL.replace(/^http/, "ws");
  return `${wsBase}/monitors/${monitorId}/ws`;
}

export function useMonitorChecks(monitorId: string) {
  const queryClient = useQueryClient();

  const query = useQuery({
    enabled: Boolean(monitorId),
    queryFn: () => listMonitorChecks(monitorId),
    queryKey: queryKeys.monitorChecks(monitorId),
  });

  useEffect(() => {
    if (!monitorId) return;

    const ws = new WebSocket(getWsUrl(monitorId));

    ws.onmessage = (event) => {
      try {
        const result = JSON.parse(event.data) as CheckResult;
        queryClient.setQueryData<CheckResult[]>(
          queryKeys.monitorChecks(monitorId),
          (oldData) => {
            const updated = [result, ...(oldData ?? [])];
            // Cap at 100 entries to prevent memory bloat
            return updated.slice(0, 100);
          }
        );
      } catch (err) {
        console.error("ws: failed to parse message:", err);
      }
    };

    ws.onerror = (err) => {
      console.error("ws: connection error:", err);
    };

    return () => {
      ws.close();
    };
  }, [monitorId, queryClient]);

  return query;
}
