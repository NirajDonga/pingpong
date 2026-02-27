export interface Metrics {
  dnsMs: number;    
  tcpMs: number;    
  tlsMs: number;    
  ttfbMs: number;   
  totalMs: number;  
}

export interface PingResult {
  sessionId: string; 
  workerId: string;  
  timestamp: string; 
  success: boolean;  
  metrics: Metrics;  
  error?: string;    
}

export interface StatusUpdate {
  status: string;
}