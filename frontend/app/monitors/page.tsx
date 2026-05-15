import { AppShell } from "@/components/layout/app-shell";
import { PageHeader } from "@/components/layout/page-header";
import { MonitorList } from "@/components/features/monitors/monitor-list";
import { ButtonLink } from "@/components/ui/button";

export default function MonitorsPage() {
  return (
    <AppShell requireAuth>
      <PageHeader
        actions={<ButtonLink href="/monitors/new">New monitor</ButtonLink>}
        description="Manage checks, intervals, ownership, and current status."
        eyebrow="Operations"
        title="Monitors"
      />
      <div className="py-8">
        <MonitorList />
      </div>
    </AppShell>
  );
}
