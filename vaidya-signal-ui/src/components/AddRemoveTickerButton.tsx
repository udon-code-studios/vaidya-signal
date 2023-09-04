export default function AddRemoveTickerButton(
  { ticker, inWatchlist }: { ticker: string; inWatchlist: boolean },
) {
  return (
    <button
      className="border-2 border-skin-accent rounded-lg p-1  text-skin-accent hover:bg-skin-accent hover:text-skin-inverted"
      title={inWatchlist ? "Remove from watchlist" : "Add to watchlist"}
      onClick={async () => {
        await fetch("/api/watchlist", {
          method: "POST",
          body: JSON.stringify({
            tickers: ticker,
            action: inWatchlist ? "remove" : "add",
          }),
        });
        window.location.reload();
      }}
    >
      {inWatchlist
        ? (
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6"
            viewBox="0 0 24 24"
          >
            <path fill="currentColor" d="M19 12.998H5v-2h14z" />
          </svg>
        )
        : (
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6"
            viewBox="0 0 24 24"
          >
            <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
          </svg>
        )}
    </button>
  );
}
