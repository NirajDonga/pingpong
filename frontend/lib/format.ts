import type { CheckResult, Incident, Monitor } from "@/lib/api/types";

export function statusOf(monitor?: Monitor) {
  if (!monitor?.enabled) {
    return "paused";
  }

  if (monitor.current_status === "up" || monitor.current_status === "down") {
    return monitor.current_status;
  }

  return "unknown";
}

export function incidentStatusOf(incident: Incident) {
  return incident.status === "resolved" ? "resolved" : "open";
}

export function checkStatusOf(check: CheckResult) {
  return check.success ? "up" : "down";
}

export function intervalLabel(seconds: number) {
  if (seconds >= 60 && seconds % 60 === 0) {
    return `${seconds / 60}m`;
  }

  return `${seconds}s`;
}

export function dateTimeLabel(value?: string) {
  if (!value) {
    return "-";
  }

  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export function shortTimeLabel(value?: string) {
  if (!value) {
    return "-";
  }

  return new Intl.DateTimeFormat(undefined, {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(new Date(value));
}
