import type { InputHTMLAttributes, ReactNode, SelectHTMLAttributes } from "react";

import { cn } from "@/lib/cn";

type FieldProps = {
  children: ReactNode;
  label: string;
};

export function Field({ children, label }: FieldProps) {
  return (
    <label className="grid gap-2 text-sm text-zinc-300">
      <span>{label}</span>
      {children}
    </label>
  );
}

export function Input({ className, ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      className={cn(
        "h-10 rounded-md border border-zinc-800 bg-black px-3 text-sm text-white outline-none transition-colors placeholder:text-zinc-600 focus:border-zinc-500",
        className,
      )}
      {...props}
    />
  );
}

export function Select({
  className,
  ...props
}: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        "h-10 rounded-md border border-zinc-800 bg-black px-3 text-sm text-white outline-none transition-colors focus:border-zinc-500",
        className,
      )}
      {...props}
    />
  );
}
