import type { ReactNode } from "react";
import Link from "next/link";

import { ButtonLink } from "@/components/ui/button";

const navItems = [
  { label: "Monitors", href: "/monitors" },
  { label: "Incidents", href: "/incidents" },
  { label: "Settings", href: "/settings" },
];

type AppShellProps = {
  children: ReactNode;
};

export function AppShell({ children }: AppShellProps) {
  return (
    <main className="min-h-screen bg-[var(--background)]">
      <header className="border-b border-[var(--border)] bg-black">
        <div className="mx-auto flex h-20 w-full max-w-6xl items-center justify-between px-6">
          <Link className="text-base font-semibold tracking-wide text-white" href="/">
            PingPong
          </Link>
          <nav className="hidden items-center gap-6 text-sm text-zinc-400 md:flex">
            {navItems.map((item) => (
              <Link
                className="transition-colors hover:text-white"
                href={item.href}
                key={item.href}
              >
                {item.label}
              </Link>
            ))}
          </nav>
          <ButtonLink href="/login" variant="secondary">
            Sign in
          </ButtonLink>
        </div>
      </header>
      <div className="mx-auto w-full max-w-6xl px-6">{children}</div>
    </main>
  );
}
