"use client";

import type { ReactNode } from "react";

import {
  EmptyState,
  ErrorState,
  LoadingState,
} from "@/components/data/query-state";
import { Panel, PanelBody, PanelHeader } from "@/components/ui/panel";
import { StatusBadge } from "@/components/ui/status-badge";
import {
  checkStatusOf,
  dateTimeLabel,
  incidentStatusOf,
  intervalLabel,
  shortTimeLabel,
  statusOf,
} from "@/lib/format";
import { useMonitorIncidents } from "@/lib/queries/use-incidents";
import { useMonitor, useMonitors } from "@/lib/queries/use-monitors";
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

  return (
    <div className="grid gap-6">
      <div className="grid gap-4 md:grid-cols-4">
        <Summary label="Status">
          <StatusBadge status={statusOf(monitor.data)} />
        </Summary>
        <Summary label="Interval">
          {intervalLabel(monitor.data.interval_seconds)}
        </Summary>
        <Summary label="Expected">{monitor.data.expected_status}</Summary>
        <Summary label="Failures">{monitor.data.consecutive_failures}</Summary>
      </div>
      <Panel>
        <PanelHeader>
          <h2 className="font-medium text-white">Recent checks</h2>
        </PanelHeader>
        <PanelBody className="overflow-x-auto">
          {checks.isLoading ? (
            <p className="text-sm text-zinc-500">Loading checks...</p>
          ) : checks.isError ? (
            <p className="text-sm text-red-300">{checks.error.message}</p>
          ) : !checks.data?.length ? (
            <p className="text-sm text-zinc-500">No checks recorded yet.</p>
          ) : (
            <table className="w-full min-w-[640px] text-left text-sm">
              <thead className="text-xs uppercase tracking-[0.16em] text-zinc-600">
                <tr>
                  <th className="pb-3 font-medium">Time</th>
                  <th className="pb-3 font-medium">Result</th>
                  <th className="pb-3 font-medium">Status</th>
                  <th className="pb-3 font-medium">Latency</th>
                  <th className="pb-3 font-medium">Worker</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--border)]">
                {checks.data.map((check) => (
                  <tr key={`${check.checkedAt}-${check.statusCode}`}>
                    <td className="py-4 text-zinc-300">
                      {shortTimeLabel(check.checkedAt)}
                    </td>
                    <td className="py-4">
                      <StatusBadge status={checkStatusOf(check)} />
                    </td>
                    <td className="py-4 text-zinc-300">{check.statusCode}</td>
                    <td className="py-4 text-zinc-300">
                      {check.responseTimeMs}ms
                    </td>
                    <td className="py-4 text-zinc-300">
                      {check.workerName || "Unknown"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </PanelBody>
      </Panel>
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
