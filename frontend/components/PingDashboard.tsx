'use client';

import { useState, useRef, useEffect, useMemo } from 'react';
import { PingResult } from '@/types/ping';
import ControlPanel from './ControlPanel';
import LatencyChart, { CHART_COLORS } from './LatencyChart';
import LogViewer from './LogViewer';
import { Activity, ArrowUp, Clock, Zap, BarChart3 } from 'lucide-react';

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

    const baseUrl = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080').replace(/\/$/, '');
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

  const stats = useMemo(() => {
    const total = results.length;
    const successful = results.filter(r => r.success);
    const rate = total === 0 ? 100 : (successful.length / total) * 100;
    const avg = successful.length === 0 ? 0 : successful.reduce((s, r) => s + r.metrics.totalMs, 0) / successful.length;
    const last = successful.length > 0 ? successful[0].metrics.totalMs : 0;
    return { total, rate, avg, last };
  }, [results]);

  const hasData = results.length > 0;

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 py-8 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-lg bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center">
            <Activity className="w-[18px] h-[18px] text-emerald-500" />
          </div>
          <div>
            <h1 className="text-lg font-semibold text-white tracking-tight">PingPong</h1>
            <p className="text-[11px] text-neutral-500">Distributed Network Monitor</p>
          </div>
        </div>
        {status === 'connected' && (
          <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
            <span className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse" />
            <span className="text-xs font-medium text-emerald-400">Live</span>
          </div>
        )}
      </div>

      {/* URL Input */}
      <ControlPanel
        targetUrl={targetUrl}
        setTargetUrl={setTargetUrl}
        status={status}
        onStart={startStream}
        onStop={stopStream}
      />

      {/* Stats */}
      {hasData && (
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          {[
            { label: 'Uptime', value: `${stats.rate.toFixed(1)}%`, icon: ArrowUp, accent: stats.rate >= 99 ? 'emerald' : stats.rate >= 95 ? 'amber' : 'red' },
            { label: 'Avg Response', value: `${stats.avg.toFixed(0)}ms`, icon: Clock, accent: 'blue' },
            { label: 'Last Response', value: `${stats.last}ms`, icon: Zap, accent: 'neutral' },
            { label: 'Total Checks', value: `${stats.total}`, icon: BarChart3, accent: 'neutral' },
          ].map((stat) => (
            <div key={stat.label} className="bg-[#141414] border border-neutral-800 rounded-xl p-4">
              <div className="flex items-center gap-2 mb-3">
                <stat.icon className={`w-4 h-4 ${
                  stat.accent === 'emerald' ? 'text-emerald-500' :
                  stat.accent === 'amber' ? 'text-amber-500' :
                  stat.accent === 'red' ? 'text-red-500' :
                  stat.accent === 'blue' ? 'text-blue-400' : 'text-neutral-500'
                }`} />
                <span className="text-[11px] text-neutral-500 uppercase tracking-wider">{stat.label}</span>
              </div>
              <p className={`text-2xl font-semibold font-mono tabular-nums ${
                stat.accent === 'emerald' ? 'text-emerald-400' :
                stat.accent === 'amber' ? 'text-amber-400' :
                stat.accent === 'red' ? 'text-red-400' :
                stat.accent === 'blue' ? 'text-blue-400' : 'text-white'
              }`}>{stat.value}</p>
            </div>
          ))}
        </div>
      )}

      {/* Response Time Chart */}
      {hasData && (
        <div className="bg-[#141414] border border-neutral-800 rounded-xl overflow-hidden">
          <div className="flex items-center justify-between px-5 py-4 border-b border-neutral-800/70">
            <h2 className="text-sm font-medium text-white">Response Time</h2>
            <div className="flex items-center gap-3">
              {uniqueWorkers.map((w, i) => (
                <div key={w} className="flex items-center gap-1.5 text-[11px] text-neutral-500">
                  <span className="w-2 h-2 rounded-full" style={{ backgroundColor: CHART_COLORS[i % CHART_COLORS.length] }} />
                  {w}
                </div>
              ))}
            </div>
          </div>
          <div className="p-2">
            <LatencyChart data={results} />
          </div>
        </div>
      )}

      {/* Empty state */}
      {!hasData && status === 'idle' && (
        <div className="py-24 text-center">
          <div className="w-14 h-14 mx-auto mb-5 bg-neutral-800/50 rounded-2xl flex items-center justify-center border border-neutral-800">
            <Activity className="w-6 h-6 text-neutral-600" />
          </div>
          <p className="text-neutral-500 text-sm mb-1">No active monitors</p>
          <p className="text-neutral-600 text-xs">Enter a URL and click Start to begin monitoring</p>
        </div>
      )}

      {/* Workers */}
      {(hasData || status !== 'idle') && (
        <div>
          <h2 className="text-sm font-medium text-white mb-3">Workers</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pb-8">
            {uniqueWorkers.length === 0 ? (
              <LogViewer workerId="Waiting for workers…" logs={[]} />
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
      )}
    </div>
  );
}