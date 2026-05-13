import { AppShell } from "@/components/layout/app-shell";
import { PageHeader } from "@/components/layout/page-header";
import { Button } from "@/components/ui/button";
import { Field, Input } from "@/components/ui/field";
import { Panel, PanelBody, PanelHeader } from "@/components/ui/panel";

export default function SettingsPage() {
  return (
    <AppShell>
      <PageHeader
        description="Keep account and API settings in one place as the product grows."
        eyebrow="Workspace"
        title="Settings"
      />
      <div className="max-w-2xl py-8">
        <Panel>
          <PanelHeader>
            <h2 className="font-medium text-white">Account</h2>
          </PanelHeader>
          <PanelBody>
            <form className="grid gap-5">
              <Field label="Email">
                <Input defaultValue="operator@pingpong.local" type="email" />
              </Field>
              <Field label="API base URL">
                <Input defaultValue="http://localhost:3001/api" />
              </Field>
              <div className="flex justify-end">
                <Button type="submit">Save settings</Button>
              </div>
            </form>
          </PanelBody>
        </Panel>
      </div>
    </AppShell>
  );
}
