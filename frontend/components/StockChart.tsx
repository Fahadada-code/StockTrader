"use client";

import {
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    Area,
    AreaChart,
} from 'recharts';
import { DailyData } from '../lib/api';

interface StockChartProps {
    history: Record<string, DailyData>;
    symbol: string;
}

const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
        return (
            <div className="glass p-4 rounded-xl border border-border/50 shadow-2xl">
                <p className="text-xs font-bold text-muted-foreground mb-1 uppercase tracking-wider">
                    {new Date(label).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })}
                </p>
                <p className="text-xl font-black text-foreground">
                    ${payload[0].value.toFixed(2)}
                </p>
            </div>
        );
    }
    return null;
};

export default function StockChart({ history, symbol }: StockChartProps) {
    const data = Object.entries(history)
        .map(([date, values]) => ({
            date,
            close: parseFloat(values.Close),
        }))
        .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

    if (data.length === 0) return null;

    const latestPrice = data[data.length - 1].close;
    const startPrice = data[0].close;
    const isPositive = latestPrice >= startPrice;

    // Using oklch colors from globals.css for JS-side chart rendering
    const color = isPositive ? '#22c55e' : '#ef4444';

    return (
        <div className="glass-dark p-6 md:p-8 rounded-3xl border border-border/30 h-[450px] relative">
            <div className="flex items-center justify-between mb-8">
                <div>
                    <h3 className="text-foreground text-xl font-bold tracking-tight">Market Analysis</h3>
                    <p className="text-muted-foreground text-sm">Historical price performance for {symbol}</p>
                </div>
                <div className="flex gap-2">
                    <div className="px-3 py-1 rounded-full glass text-[10px] font-bold uppercase tracking-wider">Daily</div>
                </div>
            </div>

            <div className="w-full h-[320px]">
                <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={data}>
                        <defs>
                            <linearGradient id="colorClose" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor={color} stopOpacity={0.2} />
                                <stop offset="95%" stopColor={color} stopOpacity={0} />
                            </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="4 4" stroke="rgba(255,255,255,0.05)" vertical={false} />
                        <XAxis
                            dataKey="date"
                            stroke="rgba(255,255,255,0.3)"
                            tick={{ fontSize: 10, fontWeight: 600 }}
                            tickLine={false}
                            axisLine={false}
                            dy={10}
                            tickFormatter={(value) => {
                                const date = new Date(value);
                                return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
                            }}
                        />
                        <YAxis
                            stroke="rgba(255,255,255,0.3)"
                            tick={{ fontSize: 10, fontWeight: 600 }}
                            tickLine={false}
                            axisLine={false}
                            dx={-10}
                            domain={['auto', 'auto']}
                            tickFormatter={(value) => `$${value}`}
                        />
                        <Tooltip content={<CustomTooltip />} cursor={{ stroke: 'rgba(255,255,255,0.1)', strokeWidth: 2 }} />
                        <Area
                            type="monotone"
                            dataKey="close"
                            stroke={color}
                            strokeWidth={3}
                            fillOpacity={1}
                            fill="url(#colorClose)"
                            animationDuration={1500}
                        />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
}
