"use client";

export default function AddTickerButton() {
  return (
    <div className="">
      <button
        onClick={async () => {
          fetch("https://vaidya-service.udon.studio/api/v1/add-ticker", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              ticker: "SPY",
            }),
          });
        }}
      >
        click me
      </button>
    </div>
  );
}
