import { Server, AlertOctagon } from 'lucide-react';
import { PingResult } from '@/types/ping';

export default function LogViewer({ logs }: { logs: PingResult[] }) {
  const getLatencyColor = (ms: number) => {
    if (ms < 100) return 'text-emerald-400 bg-emerald-400/10 border-emerald-400/20';
    if (ms < 300) return 'text-amber-400 bg-amber-400/10 border-amber-400/20';
    return 'text-rose-400 bg-rose-400/10 border-rose-400/20';
  };

  return (
    <div className="bg-white/5 border border-white/10 rounded-2xl flex flex-col backdrop-blur-sm overflow-hidden h-[500px] lg:h-auto">
      <div className="p-4 border-b border-white/10 bg-black/20">
        <h2 className="font-semibold flex items-center gap-2 text-white">
          <Server className="w-5 h-5 text-gray-400" />
          Node Responses
        </h2>
      </div>
      
      <div className="flex-1 overflow-y-auto p-4 space-y-3 custom-scrollbar">
        {logs.length === 0 ? (
          <p className="text-center text-gray-500 text-sm mt-10">No logs yet.</p>
        ) : (
          logs.map((res, i) => {
            // --- ERROR STATE UI ---
            if (!res.success) {
              return (
                <div key={`${res.workerId}-${res.timestamp}-${i}`} className="bg-rose-500/5 border border-rose-500/20 rounded-xl p-4">
                  <div className="flex justify-between items-start mb-2">
                    <div className="flex items-center gap-2">
                      <AlertOctagon className="w-4 h-4 text-rose-500" />
                      <span className="font-mono text-sm text-rose-300">{res.workerId}</span>
                    </div>
                    <span className="text-xs font-bold px-2.5 py-1 rounded-full border text-rose-400 bg-rose-400/10 border-rose-400/20 uppercase tracking-wider">
                      Failed
                    </span>
                  </div>
                  <div className="text-xs text-rose-400/80 font-mono mt-2 bg-rose-500/10 p-2.5 rounded-lg border border-rose-500/10">
                    {res.error || "Connection timeout or host unreachable"}
                  </div>
                </div>
              );
            }

            // --- SUCCESS STATE UI ---
            const total = res.metrics.totalMs;
            const statusColor = getLatencyColor(total);
            
            return (
              <div key={`${res.workerId}-${res.timestamp}-${i}`} className="bg-black/40 border border-white/5 rounded-xl p-4 hover:border-white/10 transition-colors">
                <div className="flex justify-between items-start mb-3">
                  <div className="flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />
                    <span className="font-mono text-sm text-gray-300">{res.workerId}</span>
                  </div>
                  <span className={`text-xs font-bold px-2.5 py-1 rounded-full border ${statusColor}`}>
                    {total} ms
                  </span>
                </div>

                <div className="flex h-1.5 w-full bg-gray-800 rounded-full overflow-hidden mb-2">
                  {total > 0 && (
                    <>
                      <div style={{ width: `${(res.metrics.dnsMs / total) * 100}%` }} className="bg-purple-500" title={`DNS: ${res.metrics.dnsMs}ms`} />
                      <div style={{ width: `${(res.metrics.tcpMs / total) * 100}%` }} className="bg-blue-500" title={`TCP: ${res.metrics.tcpMs}ms`} />
                      <div style={{ width: `${(res.metrics.tlsMs / total) * 100}%` }} className="bg-amber-500" title={`TLS: ${res.metrics.tlsMs}ms`} />
                      <div style={{ width: `${(res.metrics.ttfbMs / total) * 100}%` }} className="bg-emerald-500" title={`TTFB: ${res.metrics.ttfbMs}ms`} />
                    </>
                  )}
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}