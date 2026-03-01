import PingDashboard from '@/components/PingDashboard';
import { Activity } from 'lucide-react';

export default function Dashboard() {
  return (
    // We strictly use h-screen and overflow-hidden to prevent page scrolling
    <main className="h-screen w-screen overflow-hidden bg-[#0a0a0a] text-gray-100 font-sans selection:bg-blue-500/30 flex flex-col">
      <div className="fixed inset-0 z-0 pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full bg-blue-900/20 blur-[120px]" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[30%] h-[30%] rounded-full bg-indigo-900/10 blur-[100px]" />
      </div>

      {/* flex-1 and min-h-0 force this container to stay within the screen */}
      <div className="flex-1 flex flex-col max-w-7xl w-full mx-auto p-6 relative z-10 min-h-0">
        <header className="mb-6 flex-shrink-0 flex items-center justify-between border-b border-white/5 pb-4">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-blue-500/10 rounded-xl border border-blue-500/20">
              <Activity className="w-6 h-6 text-blue-400" />
            </div>
            <div>
              <h1 className="text-2xl font-bold tracking-tight text-white">
                Ping<span className="text-blue-500">Pong</span>
              </h1>
              <p className="text-xs text-gray-400 mt-1">Distributed Global Network Telemetry</p>
            </div>
          </div>
        </header>

        <PingDashboard />
      </div>
    </main>
  );
}