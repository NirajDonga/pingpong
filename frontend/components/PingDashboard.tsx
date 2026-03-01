'use client';

import { useState, useRef, useEffect, useMemo } from 'react';
import { PingResult } from '@/types/ping';
import ControlPanel from './ControlPanel';
import LatencyChart from './LatencyChart';
import LogViewer from './LogViewer';

export default function PingDashboard() {
  const [targetUrl, setTargetUrl] = useState('example.com');
  const [results, setResults] = useState<PingResult[]>([]);
  const [status, setStatus] = useState<'idle' | 'connected' | 'completed' | 'error'>('idle');
  
  const eventSourceRef = useRef<EventSource | null>(null);

  const startStream = () => {
    if (eventSourceRef.current) eventSourceRef.current.close();

    let finalUrl = targetUrl.trim();
    if (!finalUrl.startsWith('http://') && !finalUrl.startsWith('https://')) {
      finalUrl = `https://${finalUrl}`;
      setTargetUrl(finalUrl); 
    }

    setStatus('connected');
    setResults([]);

    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const url = new URL(`${baseUrl}/api/stream`);
    url.searchParams.append('target', finalUrl); 

    const eventSource = new EventSource(url.toString());
    eventSourceRef.current = eventSource;

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.status === 'completed') {
        setStatus('completed');
        eventSource.close();
        return;
      }
      setResults((prev) => {
        const newResults = [data, ...prev];
        return newResults.length > 500 ? newResults.slice(0, 500) : newResults;
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

  useEffect(() => {
    return () => eventSourceRef.current?.close();
  }, []);

  const uniqueWorkers = useMemo(() => Array.from(new Set(results.map(r => r.workerId))), [results]);

  // Stats for ControlPanel
  const { totalPings, successRate, avgLatency } = useMemo(() => {
    const total = results.length;
    const successful = results.filter(r => r.success);
    const rate = total === 0 ? 100 : (successful.length / total) * 100;
    const avg = successful.length === 0
      ? 0
      : successful.reduce((s, r) => s + r.metrics.totalMs, 0) / successful.length;
    return { totalPings: total, successRate: rate, avgLatency: avg };
  }, [results]);

  return (
    <div className="flex flex-col flex-1 min-h-0 gap-5">
      <div className="flex-shrink-0">
        <ControlPanel
          targetUrl={targetUrl}
          setTargetUrl={setTargetUrl}
          status={status}
          onStart={startStream}
          onStop={stopStream}
          totalPings={totalPings}
          successRate={successRate}
          avgLatency={avgLatency}
        />
      </div>

      {/* Chart */}
      <div className="h-[38%] min-h-[220px] flex-shrink-0">
        <LatencyChart data={results} status={status} />
      </div>

      {/* Worker Log Grid */}
      <div className="flex-1 min-h-0">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 h-full overflow-x-auto pb-2">
          {uniqueWorkers.length === 0 ? (
            <LogViewer workerId="Awaiting Node Connection…" logs={[]} />
          ) : (
            uniqueWorkers.map(worker => (
              <LogViewer
                key={worker}
                workerId={worker}
                logs={results.filter(r => r.workerId === worker)}
              />
            ))
          )}
        </div>
      </div>
    </div>
  );
}