import { Play, Square, Globe, Wifi, CheckCircle2, XCircle, Clock } from 'lucide-react';

interface ControlPanelProps {
  targetUrl: string;
  setTargetUrl: (url: string) => void;
  status: 'idle' | 'connected' | 'completed' | 'error';
  onStart: () => void;
  onStop: () => void;
  totalPings?: number;
  successRate?: number;
  avgLatency?: number;
}

const STATUS_CONFIG = {
  idle:      { label: 'Idle',  dot: 'bg-gray-500',                    border: 'border-white/10' },
  connected: { label: 'Live',  dot: 'bg-emerald-400 animate-ping',    border: 'border-emerald-500/30' },
  completed: { label: 'Done',  dot: 'bg-blue-400',                    border: 'border-blue-500/30' },
  error:     { label: 'Error', dot: 'bg-rose-500',                    border: 'border-rose-500/30' },
};

export default function ControlPanel({
  targetUrl, setTargetUrl, status, onStart, onStop,
  totalPings = 0, successRate = 100, avgLatency = 0,
}: ControlPanelProps) {
  const cfg = STATUS_CONFIG[status];

  const handleKey = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && status !== 'connected' && targetUrl) onStart();
  };

  return (
    <div className={`bg-white/5 border ${cfg.border} rounded-2xl p-2 backdrop-blur-sm transition-colors duration-500`}>
      <div className="flex flex-col sm:flex-row gap-2">
        <div className="relative flex-1 group">
          <Globe className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500 group-focus-within:text-blue-400 transition-colors" />
          <input
            type="url"
            value={targetUrl}
            onChange={(e) => setTargetUrl(e.target.value)}
            onKeyDown={handleKey}
            placeholder="example.com or https://..."
            disabled={status === 'connected'}
            className="w-full bg-black/40 border border-white/5 rounded-xl pl-12 pr-4 py-4 text-white placeholder-gray-600 focus:outline-none focus:border-blue-500/50 focus:bg-black/60 transition-all disabled:opacity-50"
          />
        </div>

        {status === 'connected' ? (
          <button
            onClick={onStop}
            className="flex items-center justify-center gap-2 bg-rose-500/10 hover:bg-rose-500/20 active:scale-95 text-rose-400 border border-rose-500/30 px-8 py-4 rounded-xl font-semibold transition-all"
          >
            <Square className="w-4 h-4" fill="currentColor" />
            Stop
          </button>
        ) : (
          <button
            onClick={onStart}
            disabled={!targetUrl}
            className="flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-500 active:scale-95 text-white px-8 py-4 rounded-xl font-semibold transition-all disabled:opacity-40 disabled:pointer-events-none shadow-lg shadow-blue-900/30"
          >
            <Play className="w-4 h-4" fill="currentColor" />
            Ping
          </button>
        )}
      </div>

      {status !== 'idle' && (
        <div className="mt-2 px-2 pb-1 flex flex-wrap items-center gap-x-6 gap-y-1 text-xs font-mono">
          <div className="flex items-center gap-2">
            <span className="relative flex h-2 w-2">
              <span className={`absolute inline-flex h-full w-full rounded-full opacity-75 ${cfg.dot}`} />
              <span className={`relative inline-flex rounded-full h-2 w-2 ${cfg.dot.replace(' animate-ping', '')}`} />
            </span>
            <span className="text-gray-400 uppercase tracking-widest">{cfg.label}</span>
          </div>
          <div className="flex items-center gap-1.5 text-gray-500">
            <Wifi className="w-3 h-3" />
            <span>{totalPings} pings</span>
          </div>
          <div className="flex items-center gap-1.5 text-emerald-400">
            <CheckCircle2 className="w-3 h-3" />
            <span>{successRate.toFixed(0)}% ok</span>
          </div>
          {avgLatency > 0 && (
            <div className="flex items-center gap-1.5 text-blue-400">
              <Clock className="w-3 h-3" />
              <span>avg {avgLatency.toFixed(0)} ms</span>
            </div>
          )}
          {successRate < 100 && (
            <div className="flex items-center gap-1.5 text-rose-400">
              <XCircle className="w-3 h-3" />
              <span>{(100 - successRate).toFixed(0)}% failed</span>
            </div>
          )}
        </div>
      )}
    </div>
  );
}