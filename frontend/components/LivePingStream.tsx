'use client';

import { useEffect, useState } from 'react';
import { PingResult } from '@/types/ping';

export default function LivePingStream() {
  const [results, setResults] = useState<PingResult[]>([]);
  const [status, setStatus] = useState<'idle' | 'connected' | 'completed' | 'error'>('idle');

  useEffect(() => {
   setStatus('connected');
    
    const baseUrl = process.env.NEXT_PUBLIC_API_URL;
    const eventSource = new EventSource(`${baseUrl}/api/stream`);

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.status === 'completed') {
        setStatus('completed');
        eventSource.close();
        return;
      }

      // OPTIMIZATION: Use functional state update and limit array to 100 items 
      // to prevent DOM bloat and memory leaks during long-running tests.
      setResults((prev) => {
        const newResults = [data, ...prev];
        return newResults.length > 100 ? newResults.slice(0, 100) : newResults;
      });
    };

    eventSource.onerror = () => {
      setStatus('error');
      eventSource.close();
    };

    return () => eventSource.close(); 
  }, []);

  return (
    <div className="grid gap-6">
      <div className="flex items-center gap-3 bg-gray-900 p-4 rounded-lg border border-gray-800">
        <StatusIndicator status={status} />
        <span className="text-sm text-gray-400 capitalize">{status}</span>
      </div>

      <div className="grid gap-3">
        {results.length === 0 && status !== 'error' ? (
          <div className="text-gray-500 text-center py-10">Waiting for stream...</div>
        ) : (
          results.map((res) => (
            <div key={`${res.sessionId}-${res.timestamp}`} className="bg-gray-800 p-4 rounded-lg border border-gray-700 flex justify-between items-center shadow-sm">
              <div className="flex items-center gap-4">
                <span className="text-blue-400 font-bold">{res.workerId}</span>
                <span className="text-xs px-2 py-1 bg-gray-900 rounded-md border border-gray-600">
                  Total: {res.metrics.totalMs}ms
                </span>
              </div>
              <div className="flex gap-3 text-xs text-gray-400">
                <span>DNS: {res.metrics.dnsMs}ms</span>
                <span>TCP: {res.metrics.tcpMs}ms</span>
                <span>TLS: {res.metrics.tlsMs}ms</span>
                <span>TTFB: {res.metrics.ttfbMs}ms</span>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

// Micro-component for the glowing dot
function StatusIndicator({ status }: { status: string }) {
  const colors = {
    idle: 'bg-gray-500',
    connected: 'bg-green-500 animate-pulse',
    completed: 'bg-blue-500',
    error: 'bg-red-500',
  };
  const colorClass = colors[status as keyof typeof colors] || colors.idle;
  return <span className={`w-3 h-3 rounded-full ${colorClass}`} />;
}