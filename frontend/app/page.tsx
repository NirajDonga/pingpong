import { AppShell } from "@/components/layout/app-shell";
import { ButtonLink } from "@/components/ui/button";

export default function Home() {
  return (
    <AppShell>
      <section className="flex min-h-[calc(100vh-81px)] items-center">
        <div className="max-w-2xl">
          <p className="text-sm font-medium uppercase tracking-[0.18em] text-[var(--muted)]">
            PingPong Web
          </p>
          <h1 className="mt-5 text-4xl font-semibold leading-tight text-white sm:text-5xl">
            Clean foundation for the monitoring dashboard.
          </h1>
          <p className="mt-5 max-w-xl text-base leading-7 text-zinc-400">
            The frontend is ready for the product screens: auth, monitors,
            checks, incidents, and status views.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <ButtonLink href="/monitors">Open monitors</ButtonLink>
            <ButtonLink href="/incidents" variant="secondary">
              View incidents
            </ButtonLink>
          </div>
        </div>
      </section>
    </AppShell>
  );
}
