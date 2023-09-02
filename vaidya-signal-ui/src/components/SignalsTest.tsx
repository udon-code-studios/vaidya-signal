import { useEffect, useState } from "react";

const fetchData = async (ticker: string) => {
  const res = await fetch(
    `https://vaidya-service.udon.studio/api/v1/get-signal-triggers?ticker=${ticker}`,
  );

  if (res.ok) {
    return await res.json();
  }
};

export default function SignalsTest() {
  const [input, setInput] = useState("");
  const [ticker, setTicker] = useState("");
  const [data, setData] = useState([]);

  useEffect(() => {
    console.log("ticker", ticker);
    fetchData(ticker).then((res) => {
      setData(res);
    });
  }, [ticker]);

  return (
    <div>
      <h1>SignalsTest</h1>
      <p>{JSON.stringify(data)}</p>
      <input
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
      />
      <button
        onClick={() => setTicker(input)}
      >
        Submit
      </button>
    </div>
  );
}
