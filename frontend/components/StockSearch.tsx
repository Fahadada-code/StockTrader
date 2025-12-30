"use client";

import { useState } from 'react';
import { Search, Loader2 } from 'lucide-react';

interface StockSearchProps {
    onSearch: (symbol: string) => void;
    loading: boolean;
}

export default function StockSearch({ onSearch, loading }: StockSearchProps) {
    const [symbol, setSymbol] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (symbol.trim()) {
            onSearch(symbol.toUpperCase());
        }
    };

    return (
        <form onSubmit={handleSubmit} className="flex gap-3 w-full group">
            <div className="relative flex-1">
                <div className="absolute inset-0 bg-primary/5 rounded-2xl blur-xl group-focus-within:bg-primary/10 transition-colors pointer-events-none" />
                <input
                    type="text"
                    value={symbol}
                    onChange={(e) => setSymbol(e.target.value)}
                    placeholder="Enter stock symbol (e.g., AAPL)"
                    className="relative w-full px-6 py-4 pl-12 glass rounded-2xl focus:ring-2 focus:ring-primary/50 focus:border-primary/50 outline-none text-foreground placeholder-muted-foreground transition-all duration-300"
                    disabled={loading}
                />
                <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground group-focus-within:text-primary transition-colors" />
            </div>
            <button
                type="submit"
                disabled={loading || !symbol.trim()}
                className="px-8 py-4 bg-primary hover:bg-primary/90 text-primary-foreground font-semibold rounded-2xl disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-primary/20 hover:shadow-primary/40 active:scale-95 flex items-center gap-2"
            >
                {loading ? (
                    <>
                        <Loader2 className="w-4 h-4 animate-spin" />
                        <span>Searching</span>
                    </>
                ) : (
                    <>
                        <Search className="w-4 h-4" />
                        <span>Search</span>
                    </>
                )}
            </button>
        </form>
    );
}
