// Matches backend shared.PingResult exactly
export interface PingResult {
  sessionId: string;
  workerName: string;
  totalMs: number;
  dnsMs: number;
  connMs: number;
  serverMs: number;
  responseMs: number;
  tlsHandshakeMs: number;
  isConnReused: boolean;
  error?: string;
}