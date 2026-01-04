"use client";

import { useState, useEffect, useRef } from 'react';
import StockSearch from '@/components/StockSearch';
import StockPrice from '@/components/StockPrice';
import StockChart from '@/components/StockChart';
import { getQuote, getHistory, EnhancedQuote, DailyData, startReplay } from '@/lib/api';
import { TrendingUp, AlertCircle, Loader2, Play, Activity } from 'lucide-react';

export interface Anomaly {
    symbol: string;
    type: string;
    confidence: number;
    details: string;
}

export default function Home() {
    const [data, setData] = useState<EnhancedQuote | null>(null);
    const [history, setHistory] = useState<Record<string, DailyData> | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [anomalies, setAnomalies] = useState<Anomaly[]>([]);

    const socketRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        const socket = new WebSocket('ws://localhost:8080/ws');
        socketRef.current = socket;

        socket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            if (msg.type === 'price') {
                setData(msg.data);
            } else if (msg.type === 'anomaly') {
                setAnomalies(prev => [msg.data, ...prev].slice(0, 5));
            }
        };

        return () => {
            socket.close();
        };
    }, []);

    const handleSearch = async (symbol: string) => {
        setLoading(true);
        setError(null);
        setData(null);
        setHistory(null);
        setAnomalies([]);

        try {
            const [quoteData, historyData] = await Promise.all([
                getQuote(symbol),
                getHistory(symbol)
            ]);

            setData(quoteData as unknown as EnhancedQuote);
            setHistory(historyData);

            // Subscribe via WebSocket
            if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
                socketRef.current.send(JSON.stringify({
                    action: 'subscribe',
                    symbol: symbol.toUpperCase()
                }));
            }
        } catch (err: any) {
            setError(err.message || 'Failed to fetch stock data. Please check the symbol and try again.');
        } finally {
            setLoading(false);
        }
    };

    const handleReplay = async () => {
        if (!data) return;
        try {
            await startReplay(data.quote.Symbol, 5.0); // 5x speed
            setAnomalies([]); // Clear old anomalies for replay
        } catch (err: any) {
            setError('Could not start replay: ' + err.message);
        }
    };

    return (
        <main className="container mx-auto px-4 py-8 md:py-16 flex flex-col items-center">
            <header className="w-full max-w-4xl text-center space-y-4 mb-12">
                <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full glass text-primary text-xs font-medium uppercase tracking-wider animate-in fade-in slide-in-from-bottom-2">
                    <TrendingUp className="w-3 h-3" />
                    Market Intelligence
                </div>
                <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight text-gradient py-2">
                    StockTrader Pro
                </h1>
                <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
                    Real-time market data and precision analytics for professional traders.
                </p>
            </header>

            <section className="w-full max-w-xl mb-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
                <StockSearch onSearch={handleSearch} loading={loading} />
            </section>

            {error && (
                <div className="w-full max-w-xl p-4 glass-dark border-destructive/30 text-destructive rounded-2xl flex items-center gap-3 animate-in zoom-in-95 duration-300 mb-6">
                    <AlertCircle className="w-5 h-5 shrink-0" />
                    <p className="text-sm font-medium">{error}</p>
                </div>
            )}

            {anomalies.length > 0 && (
                <div className="w-full max-w-xl space-y-2 mb-8 animate-in slide-in-from-top-4">
                    {anomalies.map((a, i) => (
                        <div key={i} className="p-3 glass-dark text-primary rounded-xl flex items-center justify-between gap-3 text-xs">
                            <div className="flex items-center gap-2">
                                <Activity className="w-4 h-4" />
                                <span className="font-bold uppercase">{a.type}</span>
                                <span className="text-muted-foreground">{a.details}</span>
                            </div>
                            <span className="bg-primary/20 px-2 py-0.5 rounded-full font-bold">{(a.confidence * 100).toFixed(0)}%</span>
                        </div>
                    ))}
                </div>
            )}

            {loading && !data && (
                <div className="flex flex-col items-center gap-4 mt-12 py-12 animate-in fade-in zoom-in-95">
                    <Loader2 className="w-10 h-10 text-primary animate-spin" />
                    <span className="text-muted-foreground font-medium">Analyzing market data...</span>
                </div>
            )}

            {data && history && (
                <div className="w-full max-w-6xl space-y-8 animate-in fade-in slide-in-from-bottom-8 duration-700">
                    <div className="flex justify-end mb-4">
                        <button
                            onClick={handleReplay}
                            className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass hover:bg-primary/10 text-primary text-sm font-bold transition-all active:scale-95"
                        >
                            <Play className="w-4 h-4 fill-current" />
                            Simulate Historical Replay
                        </button>
                    </div>
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                        <div className="lg:col-span-1">
                            <StockPrice data={data} />
                        </div>
                        <div className="lg:col-span-2">
                            <StockChart history={history} symbol={data.quote.Symbol} />
                        </div>
                    </div>
                </div>
            )}

            <footer className="mt-20 py-8 border-t border-border w-full max-w-4xl text-center text-muted-foreground text-sm">
                &copy; 2024 StockTrader Pro. Data provided by Alpha Vantage.
            </footer>
        </main>
    );
}
