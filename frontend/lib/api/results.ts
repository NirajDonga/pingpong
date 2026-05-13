import { api } from "@/lib/api/client";
import type { CheckResult } from "@/lib/api/types";

export function listMonitorChecks(monitorId: string, limit = 100) {
  return api<CheckResult[]>(`/monitors/${monitorId}/checks?limit=${limit}`);
}
