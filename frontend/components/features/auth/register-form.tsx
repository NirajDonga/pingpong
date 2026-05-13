"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { Button } from "@/components/ui/button";
import { Field, Input } from "@/components/ui/field";
import { useRegister } from "@/lib/queries/use-auth";

export function RegisterForm() {
  const register = useRegister();
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    register.mutate(
      { email, password },
      {
        onSuccess: () => {
          router.push("/monitors");
        },
      },
    );
  }

  return (
    <form className="mt-8 grid gap-4" onSubmit={onSubmit}>
      <Field label="Email">
        <Input
          onChange={(event) => setEmail(event.target.value)}
          placeholder="you@example.com"
          type="email"
          value={email}
        />
      </Field>
      <Field label="Password">
        <Input
          onChange={(event) => setPassword(event.target.value)}
          placeholder="Minimum 8 characters"
          type="password"
          value={password}
        />
      </Field>
      {register.error ? (
        <p className="text-sm text-red-300">{register.error.message}</p>
      ) : null}
      <Button
        className="mt-2 w-full"
        disabled={register.isPending}
        type="submit"
      >
        {register.isPending ? "Creating..." : "Create account"}
      </Button>
    </form>
  );
}
