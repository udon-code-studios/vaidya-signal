import { sortWatchlistByLastTrigger } from "./pure";

const vaidyaServiceUrl = import.meta.env.VAIDYA_SERVICE_URL;

export interface SignalTrigger {
  trigger_date: string;
  low_2_date: string;
  low_1_date: string;
}

export const getSignalTriggers = async (ticker: string) => {
  if (!ticker) {
    return;
  }

  const res = await fetch(`${vaidyaServiceUrl}/api/v1/get-signal-triggers?tickers=${ticker}`);

  if (!res.ok) {
    return;
  }

  const json = await res.json();

  return json as { [ticker: string]: SignalTrigger[] };
};

export const getMostRecentSignalTrigger = async (ticker: string) => {
  const triggers = await getSignalTriggers(ticker);

  if (!triggers) {
    return;
  }

  return triggers[ticker][triggers[ticker].length - 1].trigger_date;
};
