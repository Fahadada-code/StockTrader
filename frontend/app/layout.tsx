import type { Metadata } from 'next';
import { Inter, Outfit } from 'next/font/google';
import './globals.css';

const inter = Inter({ subsets: ['latin'], variable: '--font-sans' });
const outfit = Outfit({ subsets: ['latin'], variable: '--font-heading' });

export const metadata: Metadata = {
    title: 'StockTrader Pro | Precision Analytics',
    description: 'Professional real-time stock data and predictive analytics platform.',
};

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en" className={`${inter.variable} ${outfit.variable}`}>
            <body className="antialiased selection:bg-primary/30">
                <div className="relative min-h-screen overflow-x-hidden">
                    {/* Background decoration */}
                    <div className="pointer-events-none fixed inset-0 z-0">
                        <div className="absolute top-[-10%] left-[-10%] h-[500px] w-[500px] rounded-full bg-primary/10 blur-[120px]" />
                        <div className="absolute bottom-[-10%] right-[-10%] h-[500px] w-[500px] rounded-full bg-success/5 blur-[120px]" />
                    </div>

                    <div className="relative z-10">
                        {children}
                    </div>
                </div>
            </body>
        </html>
    );
}
