import { Play, Square, Globe } from 'lucide-react';

interface ControlPanelProps {
  targetUrl: string;
  setTargetUrl: (url: string) => void;
  status: 'idle' | 'connected' | 'completed' | 'error';
  onStart: () => void;
  onStop: () => void;
}

export default function ControlPanel({ targetUrl, setTargetUrl, status, onStart, onStop }: ControlPanelProps) {
  return (
    <div className="bg-white/5 border border-white/10 rounded-2xl p-2 backdrop-blur-sm">
      <div className="flex flex-col sm:flex-row gap-2">
        <div className="relative flex-1">
          <Globe className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
          <input 
            type="url" 
            value={targetUrl}
            onChange={(e) => setTargetUrl(e.target.value)}
            placeholder="example.com"
            disabled={status === 'connected'}
            className="w-full bg-black/40 border border-white/5 rounded-xl pl-12 pr-4 py-4 text-white focus:outline-none focus:border-blue-500/50 transition-all disabled:opacity-50"
          />
        </div>
        
        {status === 'connected' ? (
          <button 
            onClick={onStop}
            className="flex items-center justify-center gap-2 bg-rose-500/10 hover:bg-rose-500/20 text-rose-500 border border-rose-500/20 px-8 py-4 rounded-xl font-semibold transition-all"
          >
            <Square className="w-5 h-5" fill="currentColor" />
            Stop Monitor
          </button>
        ) : (
          <button 
            onClick={onStart}
            disabled={!targetUrl}
            className="flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-500 text-white px-8 py-4 rounded-xl font-semibold transition-all disabled:opacity-50"
          >
            <Play className="w-5 h-5" fill="currentColor" />
            Start Ping
          </button>
        )}
      </div>
    </div>
  );
}