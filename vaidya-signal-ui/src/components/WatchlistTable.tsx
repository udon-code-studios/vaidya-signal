import { useEffect, useState } from "react";
import {
  daysAgo,
  sortWatchlistByLastTrigger,
  sortWatchlistByTicker,
} from "../utils/pure";
import { removeFromWatchlist } from "../utils/supabase";

export default function WatchlistTable() {
  const [watchlist, setWatchlist] = useState<
    { last_trigger: string | null; ticker: string }[]
  >([]);

  // get watchlist on page load
  useEffect(() => {
    (async () => {
      const res = await fetch("/api/watchlist");
      const data = await res.json();
      setWatchlist(sortWatchlistByLastTrigger(data));
    })();
  }, []);

  return (
    <div className="w-full grid grid-cols-7 px-4 sm:px-20">
      {/* ticker column header */}
      <button
        className="col-span-3 font-bold text-2xl text-center flex justify-center items-center gap-3 hover:text-skin-accent ml-4"
        onClick={() => {
          setWatchlist(sortWatchlistByTicker(watchlist));
        }}
      >
        TICKER
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          className="w-5 h-5"
        >
          <path
            fill="currentColor"
            d="M19 17h3l-4 4l-4-4h3V3h2m-8 10v2l-3.33 4H11v2H5v-2l3.33-4H5v-2M9 3H7c-1.1 0-2 .9-2 2v6h2V9h2v2h2V5a2 2 0 0 0-2-2m0 4H7V5h2Z"
          />
        </svg>
      </button>

      {/* signal trigger column header */}
      <button
        className="col-span-3 font-bold text-2xl text-center flex justify-center items-center gap-4 hover:text-skin-accent ml-4"
        onClick={() => {
          setWatchlist(sortWatchlistByLastTrigger(watchlist));
        }}
      >
        LAST SIGNAL
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="w-5 h-5"
          viewBox="0 0 24 24"
        >
          <path
            fill="currentColor"
            d="M21 17h3l-4 4l-4-4h3V3h2v14M8 16h3v-3H8v3m5-11h-1V3h-2v2H6V3H4v2H3c-1.11 0-2 .89-2 2v11c0 1.11.89 2 2 2h10c1.11 0 2-.89 2-2V7c0-1.11-.89-2-2-2M3 18v-7h10v7H3Z"
          />
        </svg>
      </button>

      {/* empty column */}
      <div className="col-span-1">
      </div>

      {/* empty row for spacing */}
      <div className="col-span-7 h-8" />

      {/* watchlist */}
      <div className="col-span-7 grid grid-cols-7 max-h-96 overflow-y-scroll">
        {watchlist.map((ticker, i) => (
          <>
            <a
              href={`/ticker/${ticker.ticker}`}
              className="col-span-6 grid grid-cols-6 text-center hover:text-skin-accent hover:underline underline-offset-2 decoration-2 py-1"
            >
              <div className="col-span-3 font-bold">
                {ticker.ticker}
              </div>
              <div className="col-span-3">
                {ticker.last_trigger
                  ? daysAgo(new Date(ticker.last_trigger))
                  : "never"}
              </div>
            </a>
            <button
              className="col-span-1 text-red-500 hover:underline decoration-2 decoration-dashed underline-offset-2 decoration-red-500 text-center w-min mx-auto px-3 py-1"
              onClick={async () => {
                await fetch("/api/watchlist", {
                  method: "POST",
                  body: JSON.stringify({
                    tickers: ticker.ticker.toLocaleUpperCase(),
                    action: "remove",
                  }),
                });
                window.location.reload();
              }}
            >
              remove
            </button>
          </>
        ))}
      </div>
    </div>
  );
}
