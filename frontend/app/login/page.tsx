import Link from "next/link";

import { LoginForm } from "@/components/features/auth/login-form";
import { Panel, PanelBody } from "@/components/ui/panel";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center px-6">
      <Panel className="w-full max-w-md">
        <PanelBody>
          <h1 className="text-2xl font-semibold text-white">Sign in</h1>
          <p className="mt-2 text-sm text-zinc-400">
            Access your monitors, incidents, and uptime history.
          </p>
          <LoginForm />
          <p className="mt-5 text-sm text-zinc-500">
            New here?{" "}
            <Link className="text-white hover:text-zinc-300" href="/register">
              Create an account
            </Link>
          </p>
        </PanelBody>
      </Panel>
    </main>
  );
}
