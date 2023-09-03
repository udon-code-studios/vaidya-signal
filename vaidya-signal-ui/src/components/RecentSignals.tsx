export default function RecentSignals() {
  const signals = [
    { ticker: "AAPL", date: "2023-09-03" },
    { ticker: "TSLA", date: "2023-09-02" },
    { ticker: "QQQ", date: "2023-09-01" },
    { ticker: "AMZN", date: "2023-08-14" },
    { ticker: "SPY", date: "2023-08-13" },
    { ticker: "GOOG", date: "2022-01-01" },
    { ticker: "MSFT", date: "2021-01-01" },
  ];

  return (
    <div className="grid space-y-2">
      {signals.map((signal, i) => (
        <a href={`/ticker/${signal.ticker}`} key={i}>
          <div className="grid grid-cols-5 hover:text-skin-accent hover:underline underline-offset-2 decoration-2 content-end">
            <div className="mx-auto col-span-2 font-bold">{signal.ticker}</div>
            <div className="mr-auto col-span-2">
              {daysAgo(new Date(signal.date))}
            </div>
            <div className="mr-auto pl-2">➔</div>
          </div>
        </a>
      ))}
    </div>
  );
}

function daysAgo(date: Date) {
  const today = new Date();
  const diff = today.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (days === 0) return "today!";
  if (days === 1) return "yesterday";

  return `${days} days ago`;
}
