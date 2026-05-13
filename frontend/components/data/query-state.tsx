import { Panel, PanelBody } from "@/components/ui/panel";

type QueryStateProps = {
  message: string;
  title: string;
};

export function LoadingState({ message, title }: QueryStateProps) {
  return (
    <Panel>
      <PanelBody>
        <p className="text-sm font-medium text-white">{title}</p>
        <p className="mt-2 text-sm text-zinc-500">{message}</p>
      </PanelBody>
    </Panel>
  );
}

export function EmptyState({ message, title }: QueryStateProps) {
  return (
    <Panel>
      <PanelBody>
        <p className="text-sm font-medium text-white">{title}</p>
        <p className="mt-2 text-sm text-zinc-500">{message}</p>
      </PanelBody>
    </Panel>
  );
}

export function ErrorState({ message, title }: QueryStateProps) {
  return (
    <Panel>
      <PanelBody>
        <p className="text-sm font-medium text-red-300">{title}</p>
        <p className="mt-2 text-sm text-zinc-500">{message}</p>
      </PanelBody>
    </Panel>
  );
}
