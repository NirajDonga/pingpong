import { Activity, AlertCircle } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { useMemo } from 'react';
import { PingResult } from '@/types/ping';

// Custom Tooltip Component
const CustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    return (
      <div className="bg-[#18181b] border border-white/10 p-4 rounded-xl shadow-xl backdrop-blur-md">
        <p className="text-gray-400 text-xs font-mono mb-2 flex items-center gap-2">
          <span className="w-2 h-2 rounded-full bg-blue-500"></span>
          {data.worker}
        </p>
        <div className="text-2xl font-bold text-white mb-3">
          {data.total} <span className="text-sm text-gray-500 font-normal">ms</span>
        </div>
        <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs font-mono">
          <div className="text-purple-400">DNS: {data.metrics.dnsMs}ms</div>
          <div className="text-blue-400">TCP: {data.metrics.tcpMs}ms</div>
          <div className="text-amber-400">TLS: {data.metrics.tlsMs}ms</div>
          <div className="text-emerald-400">TTFB: {data.metrics.ttfbMs}ms</div>
        </div>
      </div>
    );
  }
  return null;
};

export default function LatencyChart({ data, status }: { data: PingResult[], status: string }) {
  // We include the full metrics object in the chart data so the tooltip can access it
  const chartData = useMemo(() => {
    return [...data].reverse().map((r, i) => ({
      index: i,
      total: r.metrics.totalMs,
      worker: r.workerId,
      metrics: r.metrics 
    }));
  }, [data]);

  return (
    <div className="bg-white/5 border border-white/10 rounded-2xl p-6 backdrop-blur-sm flex flex-col min-h-[350px] lg:col-span-2">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold flex items-center gap-2 text-white">
          <Activity className="w-5 h-5 text-blue-400" />
          Live Latency
        </h2>
        <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-black/40 border border-white/5">
          <span className={`w-2 h-2 rounded-full ${status === 'connected' ? 'bg-emerald-500 animate-pulse' : 'bg-gray-500'}`} />
          <span className="text-xs font-bold tracking-wider text-gray-300 uppercase">{status}</span>
        </div>
      </div>
      
      <div className="flex-1 w-full relative">
        {data.length > 0 ? (
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="colorTotal" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <XAxis dataKey="index" hide />
              <YAxis stroke="#52525b" fontSize={12} tickFormatter={(val) => `${val}ms`} width={50} />
              <Tooltip content={<CustomTooltip />} cursor={{ stroke: 'rgba(255,255,255,0.1)', strokeWidth: 1, strokeDasharray: '4 4' }} />
              <Area type="monotone" dataKey="total" stroke="#3b82f6" strokeWidth={2} fillOpacity={1} fill="url(#colorTotal)" isAnimationActive={false} />
            </AreaChart>
          </ResponsiveContainer>
        ) : (
          <div className="absolute inset-0 flex flex-col items-center justify-center text-gray-500">
            <AlertCircle className="w-8 h-8 mb-2 opacity-50" />
            <p>Awaiting telemetry data...</p>
          </div>
        )}
      </div>
    </div>
  );
}