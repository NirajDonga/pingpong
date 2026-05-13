import Link from "next/link";

import { RegisterForm } from "@/components/features/auth/register-form";
import { Panel, PanelBody } from "@/components/ui/panel";

export default function RegisterPage() {
  return (
    <main className="flex min-h-screen items-center justify-center px-6">
      <Panel className="w-full max-w-md">
        <PanelBody>
          <h1 className="text-2xl font-semibold text-white">Create account</h1>
          <p className="mt-2 text-sm text-zinc-400">
            Start tracking uptime from the first monitor.
          </p>
          <RegisterForm />
          <p className="mt-5 text-sm text-zinc-500">
            Already have an account?{" "}
            <Link className="text-white hover:text-zinc-300" href="/login">
              Sign in
            </Link>
          </p>
        </PanelBody>
      </Panel>
    </main>
  );
}
