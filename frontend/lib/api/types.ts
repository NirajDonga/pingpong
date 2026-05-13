export type User = {
  email: string;
  id: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

export type Monitor = {
  id: string;
  user_id: string;
  name: string;
  url: string;
  interval_seconds: number;
  timeout_seconds: number;
  expected_status: number;
  enabled: boolean;
  current_status: "up" | "down" | "unknown" | string;
  consecutive_failures: number;
  consecutive_successes: number;
  next_check_at: string;
  created_at: string;
};

export type CreateMonitorInput = {
  name: string;
  url: string;
  interval_seconds: number;
  timeout_seconds: number;
  expected_status: number;
};

export type CheckResult = {
  monitorId: string;
  checkedAt: string;
  success: boolean;
  statusCode: number;
  responseTimeMs: number;
  dnsMs: number;
  tcpMs: number;
  tlsMs: number;
  ttfbMs: number;
  error: string;
  workerName: string;
};

export type Incident = {
  id: string;
  monitor_id: string;
  started_at: string;
  ended_at?: string;
  status: "open" | "resolved" | string;
  reason: string;
  created_at: string;
};
