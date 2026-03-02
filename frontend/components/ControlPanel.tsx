import { Play, Square } from 'lucide-react';

interface ControlPanelProps {
  targetUrl: string;
  setTargetUrl: (url: string) => void;
  status: 'idle' | 'connected' | 'completed' | 'error';
  onStart: () => void;
  onStop: () => void;
}

export default function ControlPanel({
  targetUrl, setTargetUrl, status, onStart, onStop,
}: ControlPanelProps) {
  const isLive = status === 'connected';

  return (
    <div className="flex items-center gap-3">
      <input
        type="url"
        value={targetUrl}
        onChange={(e) => setTargetUrl(e.target.value)}
        onKeyDown={(e) => { if (e.key === 'Enter' && !isLive && targetUrl) onStart(); }}
        placeholder="Enter URL to monitor…"
        disabled={isLive}
        className="flex-1 bg-[#141414] border border-neutral-800 rounded-xl px-4 py-3 text-[15px] text-white placeholder:text-neutral-600 focus:outline-none focus:border-emerald-500/50 focus:ring-1 focus:ring-emerald-500/20 disabled:opacity-50 transition-all"
      />
      {isLive ? (
        <button
          onClick={onStop}
          className="flex items-center gap-2 h-[46px] px-5 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-xl font-medium text-sm transition-all border border-red-500/20 cursor-pointer"
        >
          <Square className="w-3.5 h-3.5" fill="currentColor" />
          Stop
        </button>
      ) : (
        <button
          onClick={onStart}
          disabled={!targetUrl}
          className="flex items-center gap-2 h-[46px] px-6 bg-emerald-500 hover:bg-emerald-400 text-black font-semibold rounded-xl text-sm transition-all disabled:opacity-30 disabled:cursor-not-allowed cursor-pointer"
        >
          <Play className="w-3.5 h-3.5" fill="currentColor" />
          Start
        </button>
      )}
    </div>
  );
}