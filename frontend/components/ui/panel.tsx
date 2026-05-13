import type { ReactNode } from "react";

import { cn } from "@/lib/cn";

type PanelProps = {
  children: ReactNode;
  className?: string;
};

export function Panel({ children, className }: PanelProps) {
  return (
    <section
      className={cn(
        "rounded-lg border border-[var(--border)] bg-[var(--panel)]",
        className,
      )}
    >
      {children}
    </section>
  );
}

export function PanelHeader({ children, className }: PanelProps) {
  return (
    <div className={cn("border-b border-[var(--border)] px-5 py-4", className)}>
      {children}
    </div>
  );
}

export function PanelBody({ children, className }: PanelProps) {
  return <div className={cn("p-5", className)}>{children}</div>;
}
