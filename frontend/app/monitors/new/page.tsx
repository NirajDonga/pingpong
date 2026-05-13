import { AppShell } from "@/components/layout/app-shell";
import { PageHeader } from "@/components/layout/page-header";
import { CreateMonitorForm } from "@/components/features/monitors/create-monitor-form";
import { Panel, PanelBody } from "@/components/ui/panel";

export default function NewMonitorPage() {
  return (
    <AppShell>
      <PageHeader
        description="Create one HTTP monitor with an interval, timeout, and expected status."
        eyebrow="Monitor setup"
        title="New monitor"
      />
      <div className="max-w-2xl py-8">
        <Panel>
          <PanelBody>
            <CreateMonitorForm />
          </PanelBody>
        </Panel>
      </div>
    </AppShell>
  );
}
