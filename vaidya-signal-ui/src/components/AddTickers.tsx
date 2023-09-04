import { useState } from "react";

export default function AddTickers() {
  const [input, setInput] = useState<string>("");

  const handleEnter = async (tickers: string) => {
    if (tickers !== "") {
      await fetch("/api/watchlist", {
        method: "POST",
        body: JSON.stringify({
          tickers: tickers.toLocaleUpperCase(),
          action: "add",
        }),
      });
      window.location.reload();
    }
  };

  return (
    <div className="w-full flex flex-col gap-2 px-4">
      {/* input bar */}
      <div className="flex gap-3">
        <input
          className="block flex-grow rounded-lg border border-skin-fill 
        border-opacity-40 bg-skin-fill py-3 pl-4
        pr-3 placeholder:italic placeholder:text-opacity-75 
        focus:border-skin-accent focus:outline-none"
          placeholder="QQQ BRK.B AAPL ..."
          type="text"
          name="search"
          value={input}
          onChange={(e) => setInput(e.currentTarget.value)}
          autoComplete="off"
          onKeyDown={(e) => {
            if (e.key === "Enter" && input !== "") {
              handleEnter(input);
            }
          }}
        />
        <button
          className="bg-skin-accent text-skin-inverted px-4 rounded-lg font-bold hover:bg-skin-fill hover:text-skin-accent border-2 border-skin-accent"
          onClick={() => handleEnter(input)}
        >
          ADD TICKERS
        </button>
      </div>

      {/* instructions */}
      <div>
        <p className="text-xs italic font-bold text-center">
          * Enter a space-separated list of tickers (e.g. "AAPL MSFT GOOG")
        </p>
      </div>
    </div>
  );
}
