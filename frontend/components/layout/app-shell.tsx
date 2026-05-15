"use client";

import type { ReactNode } from "react";
import { useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { Button, ButtonLink } from "@/components/ui/button";
import { useCurrentUser, useLogout } from "@/lib/queries/use-auth";

const navItems = [
  { label: "Monitors", href: "/monitors" },
  { label: "Incidents", href: "/incidents" },
  { label: "Settings", href: "/settings" },
];

type AppShellProps = {
  children: ReactNode;
  requireAuth?: boolean;
};

export function AppShell({ children, requireAuth = false }: AppShellProps) {
  const { data: user, isLoading } = useCurrentUser();
  const logout = useLogout();
  const router = useRouter();

  useEffect(() => {
    if (requireAuth && !isLoading && !user) {
      router.replace("/login");
    }
  }, [isLoading, requireAuth, router, user]);

  function onLogout() {
    logout.mutate(undefined, {
      onSuccess: () => {
        router.push("/login");
      },
    });
  }

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
          {isLoading ? (
            <span className="h-10 w-20 rounded-md border border-zinc-800 bg-zinc-950" />
          ) : user ? (
            <div className="flex items-center gap-3">
              <Button
                disabled={logout.isPending}
                onClick={onLogout}
                variant="secondary"
              >
                Sign out
              </Button>
            </div>
          ) : (
            <ButtonLink href="/login" variant="secondary">
              Sign in
            </ButtonLink>
          )}
        </div>
      </header>
      <div className="mx-auto w-full max-w-6xl px-6">
        {requireAuth && (isLoading || !user) ? null : children}
      </div>
    </main>
  );
}
