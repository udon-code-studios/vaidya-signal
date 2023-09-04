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

  const res = await fetch(`${vaidyaServiceUrl}/api/v1/get-signal-triggers?ticker=${ticker}`);

  if (res.ok) {
    const json = (await res.json()) as [];

    if (json.length === 0) {
      return;
    }

    return json as SignalTrigger[];
  }
};

export const getMostRecentSignalTrigger = async (ticker: string) => {
  const triggers = await getSignalTriggers(ticker);

  if (!triggers) {
    return;
  }

  return triggers[triggers.length - 1].trigger_date;
};
