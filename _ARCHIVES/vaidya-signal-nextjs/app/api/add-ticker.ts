import { NextApiRequest, NextApiResponse } from "next";
import { NextResponse } from "next/server";

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const backendResponse = await fetch("https://vaidya-signal-service-dlj42m4pha-uk.a.run.app/add-ticker", {
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      ticker: "SPY",
    }),
  });

  // const data = await res.json()
  const status = backendResponse.status;

  return res;
}
