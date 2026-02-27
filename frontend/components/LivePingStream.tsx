'use client';

import { useState, useRef, useEffect } from 'react';
import { PingResult } from '@/types/ping';

export default function LivePingStream() {
  const [targetUrl, setTargetUrl] = useState('https://example.com');
  const [results, setResults] = useState<PingResult[]>([]);
  const [status, setStatus] = useState<'idle' | 'connected' | 'completed' | 'error'>('idle');
  
  // We use a ref to keep track of the active connection so we can close it
  const eventSourceRef = useRef<EventSource | null>(null);

  const startStream = () => {
    // 1. Close any existing connection
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
    }

    // 2. Auto-prepend https:// if the user forgot it
    let finalUrl = targetUrl.trim();
    if (!finalUrl.startsWith('http://') && !finalUrl.startsWith('https://')) {
      finalUrl = `https://${finalUrl}`;
      setTargetUrl(finalUrl); // Update the input box to show the corrected URL
    }

    // 3. Reset state for new run
    setStatus('connected');
    setResults([]);

    // 4. Connect with the dynamic target URL from .env.local
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const url = new URL(`${baseUrl}/api/stream`);
    url.searchParams.append('target', finalUrl); // Pass the URL to Go

    const eventSource = new EventSource(url.toString());
    eventSourceRef.current = eventSource;

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.status === 'completed') {
        setStatus('completed');
        eventSource.close();
        return;
      }

      // OPTIMIZATION: Limit array to 100 items to prevent memory leaks
      setResults((prev) => {
        const newResults = [data, ...prev];
        return newResults.length > 100 ? newResults.slice(0, 100) : newResults;
      });
    };

    eventSource.onerror = () => {
      setStatus('error');
      eventSource.close();
    };
  };

  const stopStream = () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      setStatus('completed');
    }
  };

  // Cleanup if the user leaves the page entirely
  useEffect(() => {
    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  }, []);

  return (
    <div className="grid gap-6">
      
      {/* --- CONTROLS UI --- */}
      <div className="flex flex-col sm:flex-row gap-4 bg-gray-900 p-4 rounded-lg border border-gray-800">
        <input 
          type="url" 
          value={targetUrl}
          onChange={(e) => setTargetUrl(e.target.value)}
          placeholder="https://example.com"
          disabled={status === 'connected'}
          className="flex-1 bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-blue-500 transition-colors disabled:opacity-50"
        />
        {status === 'connected' ? (
          <button 
            onClick={stopStream}
            className="bg-red-600 hover:bg-red-700 text-white px-8 py-3 rounded-lg font-bold transition-colors"
          >
            Stop
          </button>
        ) : (
          <button 
            onClick={startStream}
            disabled={!targetUrl}
            className="bg-blue-600 hover:bg-blue-500 text-white px-8 py-3 rounded-lg font-bold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Start Ping
          </button>
        )}
      </div>

      {/* --- STATUS INDICATOR --- */}
      <div className="flex items-center gap-3 px-2">
        <StatusIndicator status={status} />
        <span className="text-sm text-gray-400 font-medium tracking-wide uppercase">
          {status}
        </span>
      </div>

      {/* --- RESULTS LIST --- */}
      <div className="grid gap-3">
        {results.length === 0 && status !== 'error' && status !== 'connected' ? (
          <div className="text-gray-500 text-center py-12 border border-dashed border-gray-800 rounded-lg">
            Enter a URL and click Start Ping to begin tracking
          </div>
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

function StatusIndicator({ status }: { status: string }) {
  const colors = {
    idle: 'bg-gray-600',
    connected: 'bg-green-500 animate-pulse shadow-[0_0_8px_rgba(34,197,94,0.6)]',
    completed: 'bg-blue-500',
    error: 'bg-red-500',
  };
  const colorClass = colors[status as keyof typeof colors] || colors.idle;
  return <span className={`w-3 h-3 rounded-full ${colorClass}`} />;
}