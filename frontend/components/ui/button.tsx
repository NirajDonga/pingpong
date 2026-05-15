import type { ButtonHTMLAttributes, ReactNode } from "react";
import Link from "next/link";

import { cn } from "@/lib/cn";

type ButtonVariant = "primary" | "secondary" | "ghost" | "danger";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant;
};

const variants: Record<ButtonVariant, string> = {
  primary: "border-white bg-white !text-black hover:bg-zinc-200",
  secondary: "border-zinc-700 bg-zinc-950 text-white hover:border-zinc-500",
  ghost: "border-transparent bg-transparent text-zinc-300 hover:text-white",
  danger: "border-red-500/40 bg-red-500/10 text-red-300 hover:bg-red-500/20 hover:border-red-500/60",
};

export function buttonClasses(variant: ButtonVariant = "primary", className?: string) {
  return cn(
    "inline-flex h-10 items-center justify-center rounded-md border px-4 text-sm font-medium transition-colors disabled:pointer-events-none disabled:opacity-50",
    variants[variant],
    className,
  );
}

export function Button({
  className,
  variant = "primary",
  type = "button",
  ...props
}: ButtonProps) {
  return (
    <button
      className={buttonClasses(variant, className)}
      type={type}
      {...props}
    />
  );
}

type ButtonLinkProps = {
  children: ReactNode;
  className?: string;
  href: string;
  variant?: ButtonVariant;
};

export function ButtonLink({
  children,
  className,
  href,
  variant = "primary",
}: ButtonLinkProps) {
  return (
    <Link className={buttonClasses(variant, className)} href={href}>
      {children}
    </Link>
  );
}
