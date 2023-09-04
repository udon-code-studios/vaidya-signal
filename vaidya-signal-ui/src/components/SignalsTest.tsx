import { useEffect, useState } from "react";
import { getSignalTriggers } from "../utils/vaidya-service";

export default function SignalsTest() {
  const [input, setInput] = useState("");
  const [ticker, setTicker] = useState("");
  const [data, setData] = useState([]);

  useEffect(() => {
    (async () => {
      const res = await getSignalTriggers(ticker);
      setData(res);
    })();
  }, [ticker]);

  return (
    <div>
      <h1>SignalsTest</h1>
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
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}
