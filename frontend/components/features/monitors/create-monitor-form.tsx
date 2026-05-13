"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { Button } from "@/components/ui/button";
import { Field, Input, Select } from "@/components/ui/field";
import { useCreateMonitor } from "@/lib/queries/use-monitors";

export function CreateMonitorForm() {
  const createMonitor = useCreateMonitor();
  const router = useRouter();
  const [name, setName] = useState("");
  const [url, setURL] = useState("");
  const [intervalSeconds, setIntervalSeconds] = useState(60);
  const [timeoutSeconds, setTimeoutSeconds] = useState(5);
  const [expectedStatus, setExpectedStatus] = useState(200);

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createMonitor.mutate(
      {
        expected_status: expectedStatus,
        interval_seconds: intervalSeconds,
        name,
        timeout_seconds: timeoutSeconds,
        url,
      },
      {
        onSuccess: (monitor) => {
          router.push(`/monitors/${monitor.id}`);
        },
      },
    );
  }

  return (
    <form className="grid gap-5" onSubmit={onSubmit}>
      <Field label="Name">
        <Input
          onChange={(event) => setName(event.target.value)}
          placeholder="API Gateway"
          value={name}
        />
      </Field>
      <Field label="URL">
        <Input
          onChange={(event) => setURL(event.target.value)}
          placeholder="https://example.com/health"
          type="url"
          value={url}
        />
      </Field>
      <div className="grid gap-5 md:grid-cols-3">
        <Field label="Interval">
          <Select
            onChange={(event) => setIntervalSeconds(Number(event.target.value))}
            value={intervalSeconds}
          >
            <option value="30">30 sec</option>
            <option value="60">1 min</option>
            <option value="300">5 min</option>
          </Select>
        </Field>
        <Field label="Timeout">
          <Select
            onChange={(event) => setTimeoutSeconds(Number(event.target.value))}
            value={timeoutSeconds}
          >
            <option value="3">3 sec</option>
            <option value="5">5 sec</option>
            <option value="10">10 sec</option>
          </Select>
        </Field>
        <Field label="Expected status">
          <Input
            inputMode="numeric"
            onChange={(event) => setExpectedStatus(Number(event.target.value))}
            value={expectedStatus}
          />
        </Field>
      </div>
      {createMonitor.error ? (
        <p className="text-sm text-red-300">{createMonitor.error.message}</p>
      ) : null}
      <div className="flex justify-end">
        <Button disabled={createMonitor.isPending} type="submit">
          {createMonitor.isPending ? "Creating..." : "Create monitor"}
        </Button>
      </div>
    </form>
  );
}
