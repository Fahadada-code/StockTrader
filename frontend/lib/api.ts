
export const API_BASE_URL = 'http://localhost:8080/api';

export interface QuoteData {
  Symbol: string;
  Open: string;
  High: string;
  Low: string;
  Price: string;
  Volume: string;
  LatestTradingDay: string;
  PreviousClose: string;
  Change: string;
  ChangePercent: string;
}

export interface DailyData {
  Open: string;
  High: string;
  Low: string;
  Close: string;
  Volume: string;
}

export async function getQuote(symbol: string): Promise<QuoteData> {
  const res = await fetch(`${API_BASE_URL}/quote?symbol=${symbol}`);
  if (res.status === 429) {
    throw new Error('API rate limit reached. Please wait a minute before searching again.');
  }
  if (!res.ok) {
    throw new Error('Failed to fetch quote');
  }
  return res.json();
}

export async function getHistory(symbol: string): Promise<Record<string, DailyData>> {
  const res = await fetch(`${API_BASE_URL}/history?symbol=${symbol}`);
  if (res.status === 429) {
    throw new Error('API rate limit reached. Please wait a minute before searching again.');
  }
  if (!res.ok) {
    throw new Error('Failed to fetch history');
  }
  return res.json();
}
