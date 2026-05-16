"use client";

import { useRouter } from "next/navigation";
import type { ReactNode } from "react";

import {
  EmptyState,
  ErrorState,
  LoadingState,
} from "@/components/data/query-state";
import { Panel, PanelBody, PanelHeader } from "@/components/ui/panel";
import { StatusBadge } from "@/components/ui/status-badge";
import { Button } from "@/components/ui/button";
import {
  checkStatusOf,
  dateTimeLabel,
  incidentStatusOf,
  intervalLabel,
  shortTimeLabel,
  statusOf,
} from "@/lib/format";
import { useMonitorIncidents } from "@/lib/queries/use-incidents";
import { useMonitor, useMonitors, useSetMonitorEnabled, useDeleteMonitor } from "@/lib/queries/use-monitors";
import { useMonitorChecks } from "@/lib/queries/use-results";

type MonitorDetailProps = {
  id: string;
};

export function MonitorDetail({ id }: MonitorDetailProps) {
  const monitor = useMonitor(id);
  const checks = useMonitorChecks(id);
  const incidents = useMonitorIncidents(id);

  if (monitor.isLoading) {
    return (
      <LoadingState
        message="Fetching monitor state from the API."
        title="Loading monitor"
      />
    );
  }

  if (monitor.isError) {
    return (
      <ErrorState message={monitor.error.message} title="Could not load monitor" />
    );
  }

  if (!monitor.data) {
    return <EmptyState message="This monitor was not found." title="Missing monitor" />;
  }

  const checksData = checks.data || [];
  
  // Calculate Global Uptime
  const totalChecks = checksData.length;
  const successfulChecks = checksData.filter((c) => c.success).length;
  const uptimePercentage =
    totalChecks > 0
      ? ((successfulChecks / totalChecks) * 100).toFixed(2) + "%"
      : "N/A";

  // Group by region
  const checksByRegion = checksData.reduce((acc, check) => {
    const region = check.workerName || "Unknown Region";
    if (!acc[region]) acc[region] = [];
    acc[region].push(check);
    return acc;
  }, {} as Record<string, typeof checksData>);

  return (
    <div className="grid gap-6">
      <div className="grid gap-4 md:grid-cols-5">
        <Summary label="Status">
          <StatusBadge status={statusOf(monitor.data)} />
        </Summary>
        <Summary label="Uptime">{uptimePercentage}</Summary>
        <Summary label="Interval">
          {intervalLabel(monitor.data.interval_seconds)}
        </Summary>
        <Summary label="Expected">{monitor.data.expected_status}</Summary>
        <Summary label="Failures">{monitor.data.consecutive_failures}</Summary>
      </div>

      <div className="space-y-6">
        <h2 className="text-xl font-medium text-white">Region Overview</h2>
        {checks.isLoading ? (
          <p className="text-sm text-zinc-500">Loading checks...</p>
        ) : checks.isError ? (
          <p className="text-sm text-red-300">{checks.error.message}</p>
        ) : Object.keys(checksByRegion).length === 0 ? (
          <p className="text-sm text-zinc-500">No checks recorded yet.</p>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {Object.entries(checksByRegion)
              .sort(([regionA], [regionB]) => regionA.localeCompare(regionB))
              .map(([region, regionChecks]) => {
              // Ensure checks are sorted by time (newest first)
              const sortedChecks = [...regionChecks].sort(
                (a, b) => new Date(b.checkedAt).getTime() - new Date(a.checkedAt).getTime()
              );
              const latestCheck = sortedChecks[0];

              const regionTotal = sortedChecks.length;
              const regionSuccess = sortedChecks.filter((c) => c.success).length;
              const regionUptime =
                regionTotal > 0
                  ? ((regionSuccess / regionTotal) * 100).toFixed(2) + "%"
                  : "N/A";
              const avgLatency =
                regionTotal > 0
                  ? (
                      sortedChecks.reduce((sum, c) => sum + c.responseTimeMs, 0) /
                      regionTotal
                    ).toFixed(0) + "ms"
                  : "N/A";

              // Take up to 20 most recent checks and reverse for chronological display L -> R
              const recentHistory = sortedChecks.slice(0, 20).reverse();

              return (
                <Panel key={region}>
                  <PanelBody className="flex flex-col gap-5">
                    <div className="flex items-center justify-between">
                      <h3 className="font-medium text-white flex items-center gap-2">
                        <span className="text-zinc-400">📍</span> {region}
                      </h3>
                      {latestCheck && (
                        <StatusBadge status={checkStatusOf(latestCheck)} />
                      )}
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.16em] text-zinc-600">
                          Uptime
                        </p>
                        <p className="mt-1 text-sm font-medium text-white">
                          {regionUptime}
                        </p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.16em] text-zinc-600">
                          Avg Latency
                        </p>
                        <p className="mt-1 text-sm font-medium text-white">
                          {avgLatency}
                        </p>
                      </div>
                    </div>

                    {recentHistory.length > 0 && (
                      <div className="mt-2">
                        <div className="flex items-center justify-between mb-2">
                          <p className="text-[10px] uppercase tracking-[0.16em] text-zinc-600">
                            Recent History
                          </p>
                          <span className="text-[10px] text-zinc-500">
                            Last {recentHistory.length}
                          </span>
                        </div>
                        <div className="flex gap-1 h-8">
                          {recentHistory.map((c, i) => (
                            <div
                              key={i}
                              title={`${shortTimeLabel(c.checkedAt)} - ${c.responseTimeMs}ms`}
                              className={`flex-1 rounded-[2px] transition-opacity hover:opacity-80 ${
                                c.success ? "bg-green-500/80" : "bg-red-500/80"
                              }`}
                            />
                          ))}
                        </div>
                      </div>
                    )}
                  </PanelBody>
                </Panel>
              );
            })}
          </div>
        )}
      </div>

      <Panel>
        <PanelHeader>
          <h2 className="font-medium text-white">Incidents</h2>
        </PanelHeader>
        <PanelBody className="grid gap-3">
          {incidents.isLoading ? (
            <p className="text-sm text-zinc-500">Loading incidents...</p>
          ) : incidents.isError ? (
            <p className="text-sm text-red-300">{incidents.error.message}</p>
          ) : !incidents.data?.length ? (
            <p className="text-sm text-zinc-500">No incidents for this monitor.</p>
          ) : (
            incidents.data.map((incident) => (
              <div
                className="flex flex-col gap-3 border-b border-[var(--border)] py-3 last:border-0 md:flex-row md:items-center md:justify-between"
                key={incident.id}
              >
                <div>
                  <p className="text-sm text-white">{incident.reason}</p>
                  <p className="mt-1 text-xs text-zinc-500">
                    {dateTimeLabel(incident.started_at)} to{" "}
                    {dateTimeLabel(incident.ended_at)}
                  </p>
                </div>
                <StatusBadge status={incidentStatusOf(incident)} />
              </div>
            ))
          )}
        </PanelBody>
      </Panel>
    </div>
  );
}

export function MonitorActions({ id }: MonitorDetailProps) {
  const router = useRouter();
  const monitor = useMonitor(id);
  const toggleEnabled = useSetMonitorEnabled(id);
  const remove = useDeleteMonitor(id);

  const isEnabled = monitor.data?.enabled ?? true;

  function handleToggle() {
    toggleEnabled.mutate(!isEnabled);
  }

  function handleDelete() {
    if (!confirm("Are you sure you want to delete this monitor? This cannot be undone.")) {
      return;
    }
    remove.mutate(undefined, {
      onSuccess: () => router.push("/monitors"),
    });
  }

  if (!monitor.data) {
    return null;
  }

  return (
    <>
      <Button
        disabled={toggleEnabled.isPending}
        onClick={handleToggle}
        variant="secondary"
        id="toggle-monitor-btn"
      >
        {toggleEnabled.isPending
          ? "Updating..."
          : isEnabled
            ? "⏸ Pause"
            : "▶ Resume"}
      </Button>
      <Button
        disabled={remove.isPending}
        onClick={handleDelete}
        variant="danger"
        id="delete-monitor-btn"
      >
        {remove.isPending ? "Deleting..." : "🗑 Delete"}
      </Button>
    </>
  );
}

export function MonitorTitle({ id }: MonitorDetailProps) {
  const monitor = useMonitor(id);
  const monitors = useMonitors();
  const fallback = monitors.data?.find((item) => item.id === id);

  return (
    <>
      <span>{monitor.data?.name ?? fallback?.name ?? "Monitor"}</span>
      <span className="mt-3 block text-sm font-normal leading-6 text-zinc-400">
        {monitor.data?.url ?? fallback?.url ?? id}
      </span>
    </>
  );
}

function Summary({ children, label }: { children: ReactNode; label: string }) {
  return (
    <Panel>
      <PanelBody>
        <p className="text-xs uppercase tracking-[0.16em] text-zinc-600">{label}</p>
        <div className="mt-3 text-sm text-white">{children}</div>
      </PanelBody>
    </Panel>
  );
}
