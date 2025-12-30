"use client";

import { QuoteData } from '../lib/api';
import { TrendingUp, TrendingDown, Activity, BarChart3 } from 'lucide-react';

interface StockPriceProps {
    quote: QuoteData;
}

export default function StockPrice({ quote }: StockPriceProps) {
    const change = parseFloat(quote.Change);
    const isPositive = change >= 0;

    return (
        <div className="glass-dark p-8 rounded-3xl border border-border/50 relative overflow-hidden group">
            {/* Background accent */}
            <div className={`absolute top-0 right-0 w-32 h-32 blur-[80px] -mr-16 -mt-16 transition-opacity duration-500 opacity-20 group-hover:opacity-40 ${isPositive ? 'bg-success' : 'bg-destructive'}`} />

            <div className="relative z-10 space-y-8">
                <div className="flex justify-between items-start">
                    <div>
                        <div className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md bg-muted text-muted-foreground text-[10px] font-bold uppercase tracking-wider mb-2">
                            Symbol
                        </div>
                        <h2 className="text-4xl font-black tracking-tight text-foreground group-hover:text-primary transition-colors duration-300">
                            {quote.Symbol}
                        </h2>
                        <p className="text-muted-foreground text-xs mt-1 flex items-center gap-1">
                            <Activity className="w-3 h-3" />
                            Updated: {quote.LatestTradingDay}
                        </p>
                    </div>
                    <div className="text-right">
                        <div className="text-4xl font-extrabold tracking-tighter text-foreground">
                            ${parseFloat(quote.Price).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                        </div>
                        <div className={`flex items-center justify-end gap-1 mt-1 font-bold ${isPositive ? 'text-success' : 'text-destructive'}`}>
                            {isPositive ? <TrendingUp className="w-4 h-4" /> : <TrendingDown className="w-4 h-4" />}
                            <span>{isPositive ? '+' : ''}{parseFloat(quote.Change).toFixed(2)}</span>
                            <span className="text-xs opacity-80">({quote.ChangePercent})</span>
                        </div>
                    </div>
                </div>

                <div className="grid grid-cols-2 gap-4 pt-6 border-t border-border/20">
                    <div className="space-y-1">
                        <span className="text-muted-foreground text-[10px] font-bold uppercase tracking-widest flex items-center gap-1">
                            <BarChart3 className="w-3 h-3" /> Volume
                        </span>
                        <div className="text-lg font-bold text-foreground">
                            {parseInt(quote.Volume).toLocaleString()}
                        </div>
                    </div>
                    <div className="space-y-1 text-right">
                        <span className="text-muted-foreground text-[10px] font-bold uppercase tracking-widest">
                            Previous Close
                        </span>
                        <div className="text-lg font-bold text-foreground">
                            ${parseFloat(quote.PreviousClose).toFixed(2)}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
