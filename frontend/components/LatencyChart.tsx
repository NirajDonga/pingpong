import { Activity, AlertCircle, TrendingDown, TrendingUp, Minus } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';
import { useMemo } from 'react';
import { PingResult } from '@/types/ping';

// A palette of distinct colors for different workers
const WORKER_COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#ef4444', '#ec4899', '#06b6d4'];

// Custom Tooltip Component
const CustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    // We grab the color assigned to the active line
    const color = payload[0].color || '#3b82f6';
    // The actual total for this specific point is stored under the worker's ID key
    const total = data[data.worker];

    return (
      <div className="bg-[#18181b] border border-white/10 p-4 rounded-xl shadow-xl backdrop-blur-md">
        <p className="text-gray-400 text-xs font-mono mb-2 flex items-center gap-2">
          <span className="w-2 h-2 rounded-full" style={{ backgroundColor: color }}></span>
          {data.worker}
        </p>
        <div className="text-2xl font-bold text-white mb-3">
          {total} <span className="text-sm text-gray-500 font-normal">ms</span>
        </div>
        {data.metrics && (
          <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs font-mono">
            <div className="text-purple-400">DNS: {data.metrics.dnsMs}ms</div>
            <div className="text-blue-400">TCP: {data.metrics.tcpMs}ms</div>
            <div className="text-amber-400">TLS: {data.metrics.tlsMs}ms</div>
            <div className="text-emerald-400">TTFB: {data.metrics.ttfbMs}ms</div>
          </div>
        )}
      </div>
    );
  }
  return null;
};

export default function LatencyChart({ data, status }: { data: PingResult[], status: string }) {
  const uniqueWorkers = useMemo(() => Array.from(new Set(data.map(d => d.workerId))), [data]);

  const chartData = useMemo(() => {
    return [...data].reverse().map((r, i) => {
      const point: any = { index: i, worker: r.workerId, metrics: r.metrics };
      point[r.workerId] = r.metrics.totalMs;
      return point;
    });
  }, [data]);

  // Compute global min / max / avg across all successful results
  const stats = useMemo(() => {
    const vals = data.filter(d => d.success).map(d => d.metrics.totalMs);
    if (vals.length === 0) return null;
    return {
      min: Math.min(...vals),
      max: Math.max(...vals),
      avg: vals.reduce((a, b) => a + b, 0) / vals.length,
    };
  }, [data]);

  return (
    <div className="bg-white/5 border border-white/10 rounded-2xl p-5 backdrop-blur-sm flex flex-col h-full w-full">
      {/* Header */}
      <div className="flex flex-wrap items-center justify-between gap-3 mb-4 flex-shrink-0">
        <h2 className="text-base font-semibold flex items-center gap-2 text-white">
          <Activity className="w-4 h-4 text-blue-400" />
          Live Latency
        </h2>

        <div className="flex flex-wrap items-center gap-3 text-xs font-mono">
          {/* Worker legend */}
          {uniqueWorkers.map((worker, i) => (
            <div key={worker} className="flex items-center gap-1.5 text-gray-400">
              <span className="w-2 h-2 rounded-full" style={{ backgroundColor: WORKER_COLORS[i % WORKER_COLORS.length] }} />
              {worker}
            </div>
          ))}

          {/* Stats pills */}
          {stats && (
            <div className="flex items-center gap-2 ml-2">
              <span className="flex items-center gap-1 px-2 py-0.5 rounded-md bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
                <TrendingDown className="w-3 h-3" /> {stats.min}ms
              </span>
              <span className="flex items-center gap-1 px-2 py-0.5 rounded-md bg-blue-500/10 border border-blue-500/20 text-blue-400">
                <Minus className="w-3 h-3" /> {stats.avg.toFixed(0)}ms
              </span>
              <span className="flex items-center gap-1 px-2 py-0.5 rounded-md bg-rose-500/10 border border-rose-500/20 text-rose-400">
                <TrendingUp className="w-3 h-3" /> {stats.max}ms
              </span>
            </div>
          )}

          {/* Status badge */}
          <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-black/40 border border-white/5">
            <span className={`w-2 h-2 rounded-full ${status === 'connected' ? 'bg-emerald-500 animate-pulse' : 'bg-gray-500'}`} />
            <span className="tracking-wider text-gray-300 uppercase">{status}</span>
          </div>
        </div>
      </div>

      {/* Chart */}
      <div className="flex-1 w-full relative min-h-0">
        {data.length > 0 ? (
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chartData} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.04)" />
              <XAxis dataKey="index" hide />
              <YAxis stroke="#52525b" fontSize={11} tickFormatter={(val) => `${val}ms`} />
              <Tooltip
                content={<CustomTooltip />}
                cursor={{ stroke: 'rgba(255,255,255,0.08)', strokeWidth: 1, strokeDasharray: '4 4' }}
              />
              {uniqueWorkers.map((workerId, index) => (
                <Line
                  key={workerId}
                  type="monotone"
                  dataKey={workerId}
                  stroke={WORKER_COLORS[index % WORKER_COLORS.length]}
                  strokeWidth={2}
                  dot={false}
                  isAnimationActive={false}
                  connectNulls
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
        ) : (
          <div className="absolute inset-0 flex flex-col items-center justify-center gap-3 text-gray-600">
            <div className="p-4 rounded-2xl bg-white/5 border border-white/5">
              <AlertCircle className="w-8 h-8 opacity-40" />
            </div>
            <p className="text-sm">Awaiting telemetry data…</p>
          </div>
        )}
      </div>
    </div>
  );
}