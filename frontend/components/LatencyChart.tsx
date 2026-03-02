'use client';

import { AreaChart, Area, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';
import { useMemo, useRef, useState, useEffect } from 'react';
import { PingResult } from '@/types/ping';

export const CHART_COLORS = ['#22c55e', '#60a5fa', '#a78bfa', '#f59e0b', '#f87171', '#ec4899'];

function ChartTooltip({ active, payload }: any) {
  if (!active || !payload?.length) return null;
  const data = payload[0]?.payload;
  return (
    <div className="bg-[#1a1a1a] border border-neutral-700/50 px-3 py-2.5 rounded-lg shadow-2xl">
      {payload.filter((p: any) => p.value != null).map((p: any) => (
        <div key={p.dataKey} className="flex items-center gap-2 py-0.5">
          <span className="w-2 h-2 rounded-full" style={{ backgroundColor: p.color }} />
          <span className="text-neutral-400 text-xs">{p.dataKey}</span>
          <span className="text-white text-xs font-mono font-semibold ml-auto pl-4">{p.value}ms</span>
        </div>
      ))}
      {data?.metrics && (
        <div className="mt-2 pt-2 border-t border-neutral-700/50 grid grid-cols-2 gap-x-4 gap-y-0.5 text-[10px] font-mono text-neutral-500">
          <span>DNS {data.metrics.dnsMs}ms</span>
          <span>TCP {data.metrics.tcpMs}ms</span>
          <span>TLS {data.metrics.tlsMs}ms</span>
          <span>TTFB {data.metrics.ttfbMs}ms</span>
        </div>
      )}
    </div>
  );
}

export default function LatencyChart({ data }: { data: PingResult[] }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [size, setSize] = useState({ width: 0, height: 0 });

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    const obs = new ResizeObserver((entries) => {
      const { width, height } = entries[0].contentRect;
      setSize({ width: Math.floor(width), height: Math.floor(height) });
    });
    obs.observe(el);
    return () => obs.disconnect();
  }, []);

  const workers = useMemo(() => Array.from(new Set(data.map(d => d.workerId))), [data]);

  const chartData = useMemo(() => {
    return [...data].reverse().map((r, i) => {
      const point: any = { index: i, worker: r.workerId, metrics: r.metrics };
      point[r.workerId] = r.success ? r.metrics.totalMs : null;
      return point;
    });
  }, [data]);

  return (
    <div ref={containerRef} className="h-[260px] w-full">
      {size.width > 0 && size.height > 0 && chartData.length > 0 && (
        <AreaChart width={size.width} height={size.height} data={chartData} margin={{ top: 12, right: 12, left: -16, bottom: 4 }}>
          <defs>
            {workers.map((w, i) => (
              <linearGradient key={w} id={`area-${i}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor={CHART_COLORS[i % CHART_COLORS.length]} stopOpacity={0.15} />
                <stop offset="100%" stopColor={CHART_COLORS[i % CHART_COLORS.length]} stopOpacity={0} />
              </linearGradient>
            ))}
          </defs>
          <CartesianGrid stroke="#1a1a1a" vertical={false} />
          <XAxis dataKey="index" hide />
          <YAxis
            tick={{ fill: '#525252', fontSize: 11 }}
            axisLine={false}
            tickLine={false}
            tickFormatter={(v) => `${v}ms`}
          />
          <Tooltip content={<ChartTooltip />} cursor={{ stroke: '#333', strokeWidth: 1 }} />
          {workers.map((w, i) => (
            <Area
              key={w}
              type="monotone"
              dataKey={w}
              stroke={CHART_COLORS[i % CHART_COLORS.length]}
              strokeWidth={1.5}
              fill={`url(#area-${i})`}
              dot={false}
              isAnimationActive={false}
              connectNulls
            />
          ))}
        </AreaChart>
      )}
    </div>
  );
}