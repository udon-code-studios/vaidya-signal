import { useEffect, useState } from "react";
import { daysAgo, sortWatchlistByLastTrigger } from "../utils/pure";

export default function RecentSignals() {
  const [watchlist, setWatchlist] = useState<
    { last_trigger: string | null; ticker: string }[]
  >([]);

  // get watchlist on page load
  useEffect(() => {
    (async () => {
      const res = await fetch("/api/watchlist");
      const data = await res.json();
      setWatchlist(sortWatchlistByLastTrigger(data).slice(0, 7)); // only show the first 7
    })();
  }, []);

  return (
    <div className="grid space-y-2">
      {watchlist.map((ticker, i) => (
        <a href={`/ticker/${ticker.ticker}`} key={i}>
          <div className="grid grid-cols-5 hover:text-skin-accent hover:underline underline-offset-2 decoration-2 content-end">
            <div className="mx-auto col-span-2 font-bold">{ticker.ticker}</div>
            <div className="mr-auto col-span-2">
              {ticker.last_trigger ? daysAgo(ticker.last_trigger) : "never"}
            </div>
            <div className="mr-auto pl-2">➔</div>
          </div>
        </a>
      ))}
    </div>
  );
}
