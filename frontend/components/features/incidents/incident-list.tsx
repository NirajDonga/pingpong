"use client";

import {
  EmptyState,
  ErrorState,
  LoadingState,
} from "@/components/data/query-state";
import { Panel, PanelBody } from "@/components/ui/panel";
import { StatusBadge } from "@/components/ui/status-badge";
import { dateTimeLabel, incidentStatusOf } from "@/lib/format";
import { useIncidents } from "@/lib/queries/use-incidents";

export function IncidentList() {
  const incidents = useIncidents();

  if (incidents.isLoading) {
    return (
      <LoadingState
        message="Fetching incidents from the API."
        title="Loading incidents"
      />
    );
  }

  if (incidents.isError) {
    return (
      <ErrorState message={incidents.error.message} title="Could not load incidents" />
    );
  }

  if (!incidents.data?.length) {
    return (
      <EmptyState
        message="Incidents will appear when a monitor enters downtime."
        title="No incidents yet"
      />
    );
  }

  return (
    <div className="grid gap-4">
      {incidents.data.map((incident) => (
        <Panel key={incident.id}>
          <PanelBody className="grid gap-4 md:grid-cols-[1fr_auto] md:items-center">
            <div>
              <div className="flex flex-wrap items-center gap-3">
                <h2 className="text-lg font-medium text-white">
                  {incident.monitor_id}
                </h2>
                <StatusBadge status={incidentStatusOf(incident)} />
              </div>
              <p className="mt-2 text-sm text-zinc-500">{incident.reason}</p>
            </div>
            <p className="text-sm text-zinc-400">
              {dateTimeLabel(incident.started_at)} to {dateTimeLabel(incident.ended_at)}
            </p>
          </PanelBody>
        </Panel>
      ))}
    </div>
  );
}
