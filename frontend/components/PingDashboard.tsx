'use client';

import { useState, useRef, useEffect } from 'react';
import { PingResult } from '@/types/ping';
import ControlPanel from './ControlPanel';
import LatencyChart from './LatencyChart';
import LogViewer from './LogViewer';

export default function PingDashboard() {
  const [targetUrl, setTargetUrl] = useState('example.com');
  const [results, setResults] = useState<PingResult[]>([]);
  const [status, setStatus] = useState<'idle' | 'connected' | 'completed' | 'error'>('idle');
  
  const eventSourceRef = useRef<EventSource | null>(null);

  const startStream = () => {
    if (eventSourceRef.current) eventSourceRef.current.close();

    let finalUrl = targetUrl.trim();
    if (!finalUrl.startsWith('http://') && !finalUrl.startsWith('https://')) {
      finalUrl = `https://${finalUrl}`;
      setTargetUrl(finalUrl); 
    }

    setStatus('connected');
    setResults([]);

    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const url = new URL(`${baseUrl}/api/stream`);
    url.searchParams.append('target', finalUrl); 

    const eventSource = new EventSource(url.toString());
    eventSourceRef.current = eventSource;

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.status === 'completed') {
        setStatus('completed');
        eventSource.close();
        return;
      }
      setResults((prev) => {
        const newResults = [data, ...prev];
        return newResults.length > 500 ? newResults.slice(0, 500) : newResults;
      });
    };

    eventSource.onerror = () => {
      setStatus('error');
      eventSource.close();
    };
  };

  const stopStream = () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      setStatus('completed');
    }
  };

  useEffect(() => {
    return () => eventSourceRef.current?.close();
  }, []);

  return (
    <div className="grid gap-6">
      <ControlPanel 
        targetUrl={targetUrl} 
        setTargetUrl={setTargetUrl} 
        status={status} 
        onStart={startStream} 
        onStop={stopStream} 
      />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <LatencyChart data={results} status={status} />
        <LogViewer logs={results} />
      </div>
    </div>
  );
}