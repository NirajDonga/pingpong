import { MonitorHistory } from "@/components/features/monitors/monitor-history";
import { MonitorTitle } from "@/components/features/monitors/monitor-detail";
import { AppShell } from "@/components/layout/app-shell";
import { PageHeader } from "@/components/layout/page-header";
import { ButtonLink } from "@/components/ui/button";

type MonitorHistoryPageProps = {
  params: Promise<{ id: string }>;
};

export default async function MonitorHistoryPage({
  params,
}: MonitorHistoryPageProps) {
  const { id } = await params;

  return (
    <AppShell requireAuth>
      <PageHeader
        actions={
          <ButtonLink href={`/monitors/${id}`} variant="secondary">
            Back to Overview
          </ButtonLink>
        }
        description="Detailed chronological log of all uptime checks and latency metrics."
        eyebrow="History"
        title={<MonitorTitle id={id} />}
      />
      <div className="py-8">
        <MonitorHistory id={id} />
      </div>
    </AppShell>
  );
}
