import { MonitorDetail, MonitorTitle } from "@/components/features/monitors/monitor-detail";
import { AppShell } from "@/components/layout/app-shell";
import { PageHeader } from "@/components/layout/page-header";
import { ButtonLink } from "@/components/ui/button";

type MonitorDetailPageProps = {
  params: Promise<{ id: string }>;
};

export default async function MonitorDetailPage({
  params,
}: MonitorDetailPageProps) {
  const { id } = await params;

  return (
    <AppShell>
      <PageHeader
        actions={
          <ButtonLink href="/monitors" variant="secondary">
            Back
          </ButtonLink>
        }
        eyebrow="Monitor detail"
        title={<MonitorTitle id={id} />}
      />
      <div className="py-8">
        <MonitorDetail id={id} />
      </div>
    </AppShell>
  );
}
