import { PingResult } from '@/types/ping';
import { useMemo } from 'react';

export default function LogViewer({ workerId, logs }: { workerId: string; logs: PingResult[] }) {
  const successCount = useMemo(() => logs.filter(l => l.success).length, [logs]);
  const uptime = logs.length > 0 ? ((successCount / logs.length) * 100).toFixed(1) : '—';
  const avgMs = useMemo(() => {
    const s = logs.filter(l => l.success);
    return s.length > 0 ? (s.reduce((a, r) => a + r.metrics.totalMs, 0) / s.length).toFixed(0) : '—';
  }, [logs]);

  return (
    <div className="bg-[#141414] border border-neutral-800 rounded-xl overflow-hidden">
      {/* Header */}
      <div className="px-4 py-3.5 border-b border-neutral-800/70">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <span className={`w-2.5 h-2.5 rounded-full ${
              logs.length === 0 ? 'bg-neutral-700' :
              successCount > 0 ? 'bg-emerald-500 shadow-[0_0_8px_rgba(34,197,94,0.4)]' :
              'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.4)]'
            }`} />
            <span className="text-sm font-medium text-white truncate max-w-[200px]">{workerId}</span>
          </div>
          {logs.length > 0 && (
            <div className="flex items-center gap-3 text-xs text-neutral-500">
              <span className={successCount === logs.length ? 'text-emerald-400' : 'text-neutral-400'}>{uptime}%</span>
              <span>{avgMs}ms</span>
              <span>{logs.length} checks</span>
            </div>
          )}
        </div>
      </div>

      {/* Uptime bars */}
      {logs.length > 0 && (
        <div className="px-4 py-3 border-b border-neutral-800/70">
          <div className="flex gap-[2px] h-7 items-end">
            {[...logs].reverse().slice(-40).map((res, i) => (
              <div
                key={i}
                className={`flex-1 rounded-[2px] min-w-[3px] transition-all ${
                  res.success ? 'bg-emerald-500/60 hover:bg-emerald-500' : 'bg-red-500/60 hover:bg-red-500'
                }`}
                style={{
                  height: res.success
                    ? `${Math.max(20, Math.min(100, (res.metrics.totalMs / 300) * 100))}%`
                    : '100%'
                }}
                title={res.success ? `${res.metrics.totalMs}ms` : res.error || 'Failed'}
              />
            ))}
          </div>
        </div>
      )}

      {/* Logs */}
      <div className="max-h-[250px] overflow-y-auto custom-scrollbar">
        {logs.length === 0 ? (
          <div className="py-12 text-center text-neutral-600 text-sm">
            Waiting for data…
          </div>
        ) : (
          logs.slice(0, 50).map((res, i) => (
            <div
              key={`${res.workerId}-${res.timestamp}-${i}`}
              className="px-4 py-2.5 flex items-center justify-between text-xs border-b border-neutral-800/40 last:border-b-0 hover:bg-white/[0.02] transition-colors"
            >
              <div className="flex items-center gap-2.5">
                <span className={`w-1.5 h-1.5 rounded-full ${res.success ? 'bg-emerald-500' : 'bg-red-500'}`} />
                <span className="text-neutral-500 font-mono tabular-nums">
                  {new Date(res.timestamp).toLocaleTimeString()}
                </span>
              </div>
              {res.success ? (
                <span className={`font-mono font-medium tabular-nums ${
                  res.metrics.totalMs < 100 ? 'text-emerald-400' :
                  res.metrics.totalMs < 300 ? 'text-amber-400' : 'text-red-400'
                }`}>
                  {res.metrics.totalMs}ms
                </span>
              ) : (
                <span className="text-red-400 truncate max-w-[180px]">{res.error || 'Failed'}</span>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}