import { api } from "@/lib/api/client";
import type { Incident } from "@/lib/api/types";

export function listIncidents(limit = 100) {
  return api<Incident[]>(`/incidents?limit=${limit}`);
}

export function listMonitorIncidents(monitorId: string, limit = 100) {
  return api<Incident[]>(`/monitors/${monitorId}/incidents?limit=${limit}`);
}
