"use client";

import Link from "next/link";

import {
  EmptyState,
  ErrorState,
  LoadingState,
} from "@/components/data/query-state";
import { Panel, PanelBody } from "@/components/ui/panel";
import { StatusBadge } from "@/components/ui/status-badge";
import { intervalLabel, statusOf } from "@/lib/format";
import { useMonitors } from "@/lib/queries/use-monitors";

export function MonitorList() {
  const monitors = useMonitors();

  if (monitors.isLoading) {
    return (
      <LoadingState
        message="Fetching monitors from the API."
        title="Loading monitors"
      />
    );
  }

  if (monitors.isError) {
    return (
      <ErrorState message={monitors.error.message} title="Could not load monitors" />
    );
  }

  if (!monitors.data?.length) {
    return (
      <EmptyState
        message="Create your first monitor to start collecting uptime data."
        title="No monitors yet"
      />
    );
  }

  return (
    <div className="grid gap-4">
      {monitors.data.map((monitor) => (
        <Panel key={monitor.id}>
          <PanelBody className="grid gap-5 md:grid-cols-[1fr_auto] md:items-center">
            <div>
              <div className="flex flex-wrap items-center gap-3">
                <Link
                  className="text-lg font-medium text-white hover:text-zinc-300"
                  href={`/monitors/${monitor.id}`}
                >
                  {monitor.name}
                </Link>
                <StatusBadge status={statusOf(monitor)} />
              </div>
              <p className="mt-2 text-sm text-zinc-500">{monitor.url}</p>
            </div>
            <div className="grid grid-cols-3 gap-6 text-sm">
              <Metric label="Interval" value={intervalLabel(monitor.interval_seconds)} />
              <Metric label="Timeout" value={`${monitor.timeout_seconds}s`} />
              <Metric label="Next check" value={new Date(monitor.next_check_at).toLocaleTimeString()} />
            </div>
          </PanelBody>
        </Panel>
      ))}
    </div>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-xs uppercase tracking-[0.16em] text-zinc-600">{label}</p>
      <p className="mt-1 text-white">{value}</p>
    </div>
  );
}
