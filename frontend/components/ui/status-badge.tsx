import { cn } from "@/lib/cn";

type Status = "up" | "down" | "unknown" | "open" | "resolved" | "paused";

const styles: Record<Status, string> = {
  up: "border-green-500/40 bg-green-500/10 text-green-300",
  down: "border-red-500/40 bg-red-500/10 text-red-300",
  unknown: "border-zinc-700 bg-zinc-900 text-zinc-300",
  open: "border-blue-500/40 bg-blue-500/10 text-blue-300",
  resolved: "border-green-500/40 bg-green-500/10 text-green-300",
  paused: "border-zinc-700 bg-zinc-900 text-zinc-300",
};

type StatusBadgeProps = {
  status: Status | string;
};

export function StatusBadge({ status }: StatusBadgeProps) {
  const knownStatus = isKnownStatus(status) ? status : "unknown";

  return (
    <span
      className={cn(
        "inline-flex h-6 items-center rounded-full border px-2.5 text-xs font-medium capitalize",
        styles[knownStatus],
      )}
    >
      {knownStatus}
    </span>
  );
}

function isKnownStatus(status: string): status is Status {
  return status in styles;
}
