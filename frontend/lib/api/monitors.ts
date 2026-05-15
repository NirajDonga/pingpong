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

export function setMonitorEnabled(id: string, enabled: boolean) {
  return api<Monitor>(`/monitors/${id}/enabled`, {
    body: { enabled },
    method: "PATCH",
  });
}

export function deleteMonitor(id: string) {
  return api<void>(`/monitors/${id}`, {
    method: "DELETE",
  });
}
