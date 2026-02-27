import LivePingStream from '@/components/LivePingStream';

export default function Dashboard() {
  return (
    <main className="min-h-screen bg-black text-white p-8 font-mono">
      <div className="max-w-5xl mx-auto">

        <header className="mb-10 border-b border-gray-800 pb-6">
          <h1 className="text-4xl font-extrabold tracking-tight mb-2">
            Ping<span className="text-blue-500">Pong</span>
          </h1>
          <p className="text-gray-400">Distributed Global Network Monitor</p>
        </header>

        <LivePingStream />

      </div>
    </main>
  );
}