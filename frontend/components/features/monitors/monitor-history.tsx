"use client";

import { useSearchParams } from "next/navigation";

import {
  EmptyState,
  ErrorState,
  LoadingState,
} from "@/components/data/query-state";
import { Panel, PanelBody } from "@/components/ui/panel";
import { StatusBadge } from "@/components/ui/status-badge";
import {
  checkStatusOf,
  dateTimeLabel,
} from "@/lib/format";
import { useMonitorChecks } from "@/lib/queries/use-results";

type MonitorHistoryProps = {
  id: string;
};

export function MonitorHistory({ id }: MonitorHistoryProps) {
  const checks = useMonitorChecks(id);

  if (checks.isLoading) {
    return (
      <LoadingState
        message="Fetching past checks from the API."
        title="Loading history"
      />
    );
  }

  if (checks.isError) {
    return (
      <ErrorState message={checks.error.message} title="Could not load history" />
    );
  }

  if (!checks.data?.length) {
    return <EmptyState message="No past checks recorded yet." title="No history" />;
  }

  const searchParams = useSearchParams();
  const regionFilter = searchParams.get("region");

  // Filter by region if specified
  const filteredChecks = regionFilter
    ? checks.data.filter((c) => (c.workerName || "Unknown") === regionFilter)
    : checks.data;

  // Sort checks by time, newest first
  const sortedChecks = [...filteredChecks].sort(
    (a, b) => new Date(b.checkedAt).getTime() - new Date(a.checkedAt).getTime()
  );

  return (
    <Panel>
      <PanelBody className="overflow-x-auto p-0">
        <table className="w-full min-w-[900px] text-left text-sm">
          <thead className="border-b border-[var(--border)] bg-black/20 text-xs uppercase tracking-[0.16em] text-zinc-600">
            <tr>
              <th className="px-5 py-4 font-medium">Time</th>
              <th className="px-5 py-4 font-medium">Region</th>
              <th className="px-5 py-4 font-medium">Result</th>
              <th className="px-5 py-4 font-medium">Code</th>
              <th className="px-5 py-4 font-medium">Total Latency</th>
              <th className="px-5 py-4 font-medium">Network Breakdown</th>
              <th className="px-5 py-4 font-medium">Error Details</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--border)]">
            {sortedChecks.map((check) => (
              <tr key={`${check.checkedAt}-${check.workerName}`} className="hover:bg-white/[0.02] transition-colors">
                <td className="px-5 py-4 text-zinc-300">
                  {dateTimeLabel(check.checkedAt)}
                </td>
                <td className="px-5 py-4 text-zinc-300">
                  {check.workerName || "Unknown"}
                </td>
                <td className="px-5 py-4">
                  <StatusBadge status={checkStatusOf(check)} />
                </td>
                <td className="px-5 py-4 text-zinc-300">{check.statusCode}</td>
                <td className="px-5 py-4 font-medium text-white">
                  {check.responseTimeMs}ms
                </td>
                <td className="px-5 py-4 text-zinc-500 text-xs flex gap-3">
                  <span title="DNS">DNS: {check.dnsMs}</span>
                  <span title="TCP">TCP: {check.tcpMs}</span>
                  <span title="TLS">TLS: {check.tlsMs}</span>
                  <span title="TTFB">TTFB: {check.ttfbMs}</span>
                </td>
                <td className="px-5 py-4 text-red-400 text-xs max-w-[200px] truncate" title={check.error}>
                  {check.error || "-"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </PanelBody>
    </Panel>
  );
}
