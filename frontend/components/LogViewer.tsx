import { Server, AlertOctagon } from 'lucide-react';
import { PingResult } from '@/types/ping';
import { useMemo } from 'react';

const METRIC_SEGMENTS = [
  { key: 'dnsMs',  label: 'DNS',  color: 'bg-purple-500',  text: 'text-purple-400' },
  { key: 'tcpMs',  label: 'TCP',  color: 'bg-blue-500',    text: 'text-blue-400' },
  { key: 'tlsMs',  label: 'TLS',  color: 'bg-amber-500',   text: 'text-amber-400' },
  { key: 'ttfbMs', label: 'TTFB', color: 'bg-emerald-500', text: 'text-emerald-400' },
] as const;

export default function LogViewer({ workerId, logs }: { workerId: string; logs: PingResult[] }) {
  const getLatencyColor = (ms: number) => {
    if (ms < 100) return 'text-emerald-400 bg-emerald-400/10 border-emerald-400/20';
    if (ms < 300) return 'text-amber-400 bg-amber-400/10 border-amber-400/20';
    return 'text-rose-400 bg-rose-400/10 border-rose-400/20';
  };

  const successLogs = useMemo(() => logs.filter(l => l.success), [logs]);
  const headerStats = useMemo(() => {
    if (successLogs.length === 0) return null;
    const vals = successLogs.map(l => l.metrics.totalMs);
    return {
      min: Math.min(...vals),
      max: Math.max(...vals),
      avg: (vals.reduce((a, b) => a + b, 0) / vals.length).toFixed(0),
    };
  }, [successLogs]);

  return (
    <div className="bg-white/5 border border-white/10 rounded-2xl flex flex-col backdrop-blur-sm overflow-hidden h-full min-h-[300px]">
      {/* Header */}
      <div className="px-4 py-3 border-b border-white/10 bg-black/20">
        <div className="flex justify-between items-center">
          <h2 className="font-semibold flex items-center gap-2 text-white text-sm">
            <Server className="w-4 h-4 text-gray-400" />
            <span className="font-mono truncate max-w-[140px]" title={workerId}>{workerId}</span>
          </h2>
          <div className="flex items-center gap-2">
            <span className="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse" />
            <span className="text-xs text-gray-500 font-mono">{logs.length}</span>
          </div>
        </div>
        {headerStats && (
          <div className="flex items-center gap-3 mt-2 text-[10px] font-mono">
            <span className="text-emerald-400">min {headerStats.min}ms</span>
            <span className="text-blue-400">avg {headerStats.avg}ms</span>
            <span className="text-rose-400">max {headerStats.max}ms</span>
          </div>
        )}
      </div>

      {/* Log list */}
      <div className="flex-1 overflow-y-auto p-3 space-y-2 custom-scrollbar">
        {logs.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full gap-2 text-gray-600 py-10">
            <Server className="w-6 h-6 opacity-30" />
            <p className="text-xs">Awaiting data…</p>
          </div>
        ) : (
          logs.map((res, i) => {
            if (!res.success) {
              return (
                <div key={`${res.workerId}-${res.timestamp}-${i}`} className="bg-rose-500/5 border border-rose-500/20 rounded-xl p-3">
                  <div className="flex items-center gap-2 mb-1.5">
                    <AlertOctagon className="w-3.5 h-3.5 text-rose-500" />
                    <span className="font-mono text-xs text-rose-300">Failed</span>
                    <span className="ml-auto text-[10px] text-gray-600 font-mono">{new Date(res.timestamp).toLocaleTimeString()}</span>
                  </div>
                  <div className="text-xs text-rose-400/80 font-mono bg-rose-500/10 p-2 rounded-lg border border-rose-500/10">
                    {res.error || 'Connection timeout'}
                  </div>
                </div>
              );
            }

            const total = res.metrics.totalMs;
            const statusColor = getLatencyColor(total);

            return (
              <div key={`${res.workerId}-${res.timestamp}-${i}`} className="bg-black/40 border border-white/5 rounded-xl p-3 hover:border-white/10 transition-colors">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-[10px] text-gray-600 font-mono">
                    {new Date(res.timestamp).toLocaleTimeString()}
                  </span>
                  <span className={`text-xs font-bold px-2 py-0.5 rounded-full border ${statusColor}`}>
                    {total} ms
                  </span>
                </div>

                {/* Stacked bar */}
                <div className="flex h-1.5 w-full bg-gray-800/80 rounded-full overflow-hidden">
                  {total > 0 && METRIC_SEGMENTS.map(seg => (
                    <div
                      key={seg.key}
                      style={{ width: `${(res.metrics[seg.key] / total) * 100}%` }}
                      className={seg.color}
                      title={`${seg.label}: ${res.metrics[seg.key]}ms`}
                    />
                  ))}
                </div>

                {/* Labels */}
                <div className="flex items-center gap-3 mt-1.5 text-[10px] font-mono">
                  {METRIC_SEGMENTS.map(seg => (
                    <span key={seg.key} className={`${seg.text} opacity-70`}>
                      {seg.label}&nbsp;{res.metrics[seg.key]}
                    </span>
                  ))}
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}