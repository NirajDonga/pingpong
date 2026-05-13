export const queryKeys = {
  incidents: ["incidents"] as const,
  monitor: (id: string) => ["monitor", id] as const,
  monitorChecks: (id: string) => ["monitor", id, "checks"] as const,
  monitorIncidents: (id: string) => ["monitor", id, "incidents"] as const,
  monitors: ["monitors"] as const,
};
