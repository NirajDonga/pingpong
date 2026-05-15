import { AppShell } from "@/components/layout/app-shell";
import { PageHeader } from "@/components/layout/page-header";
import { IncidentList } from "@/components/features/incidents/incident-list";

export default function IncidentsPage() {
  return (
    <AppShell requireAuth>
      <PageHeader
        description="Review open downtime windows and resolved recovery periods."
        eyebrow="Reliability"
        title="Incidents"
      />
      <div className="py-8">
        <IncidentList />
      </div>
    </AppShell>
  );
}
