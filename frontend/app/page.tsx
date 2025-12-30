"use client";

import { useState } from 'react';
import StockSearch from '@/components/StockSearch';
import StockPrice from '@/components/StockPrice';
import StockChart from '@/components/StockChart';
import { getQuote, getHistory, QuoteData, DailyData } from '@/lib/api';
import { TrendingUp, AlertCircle, Loader2 } from 'lucide-react';

export default function Home() {
    const [quote, setQuote] = useState<QuoteData | null>(null);
    const [history, setHistory] = useState<Record<string, DailyData> | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSearch = async (symbol: string) => {
        setLoading(true);
        setError(null);
        setQuote(null);
        setHistory(null);

        try {
            const [quoteData, historyData] = await Promise.all([
                getQuote(symbol),
                getHistory(symbol)
            ]);

            setQuote(quoteData);
            setHistory(historyData);
        } catch (err: any) {
            setError(err.message || 'Failed to fetch stock data. Please check the symbol and try again.');
        } finally {
            setLoading(false);
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
                <div className="w-full max-w-md p-4 glass-dark border-destructive/30 text-destructive rounded-2xl flex items-center gap-3 animate-in zoom-in-95 duration-300">
                    <AlertCircle className="w-5 h-5 shrink-0" />
                    <p className="text-sm font-medium">{error}</p>
                </div>
            )}

            {loading && !quote && (
                <div className="flex flex-col items-center gap-4 mt-12 py-12 animate-in fade-in zoom-in-95">
                    <Loader2 className="w-10 h-10 text-primary animate-spin" />
                    <span className="text-muted-foreground font-medium">Analyzing market data...</span>
                </div>
            )}

            {quote && history && (
                <div className="w-full max-w-6xl space-y-8 animate-in fade-in slide-in-from-bottom-8 duration-700">
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                        <div className="lg:col-span-1">
                            <StockPrice quote={quote} />
                        </div>
                        <div className="lg:col-span-2">
                            <StockChart history={history} symbol={quote.Symbol} />
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
