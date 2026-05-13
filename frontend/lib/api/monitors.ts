import { api } from "@/lib/api/client";
import type { CreateMonitorInput, Monitor } from "@/lib/api/types";

export function listMonitors() {
  return api<Monitor[]>("/monitors");
}

export function getMonitor(id: string) {
  return api<Monitor>(`/monitors/${id}`);
}

export function createMonitor(input: CreateMonitorInput) {
  return api<Monitor>("/monitors", {
    body: input,
    method: "POST",
  });
}
