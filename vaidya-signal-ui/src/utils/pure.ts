export function daysAgo(date: Date) {
  const today = new Date();
  const diff = today.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (days === 0) return "today!";
  if (days === 1) return "yesterday";

  return `${days} days ago`;
}

export const sortWatchlistByTicker = (list: { last_trigger: string | null; ticker: string }[]) => {
  const sorted = [...list].sort((a, b) => {
    if (a.ticker < b.ticker) return -1;
    if (a.ticker > b.ticker) return 1;
    return 0;
  });
  return sorted;
};

export const sortWatchlistByLastTrigger = (list: { last_trigger: string | null; ticker: string }[]) => {
  const sorted = [...list].sort((a, b) => {
    if (a.last_trigger === null) return 1;
    if (b.last_trigger === null) return -1;
    if (a.last_trigger < b.last_trigger) return 1;
    if (a.last_trigger > b.last_trigger) return -1;
    return 0;
  });
  return sorted;
};
