import type { APIRoute } from "astro";
import { addToWatchlist, getWatchlist, removeFromWatchlist } from "../../utils/supabase";

export const GET: APIRoute = async ({ request }) => {
  // get watchlist from supabase
  const watchlist = await getWatchlist();

  // return watchlist
  return new Response(JSON.stringify(watchlist), {
    status: 200,
    headers: {
      "Content-Type": "application/json",
    },
  });
};

// add or remove tickers from watchlist
// request body: { tickers: string, action: "add" | "remove" }
export const POST: APIRoute = async ({ request }) => {
  // get ticker and action from request body
  const { tickers, action } = await request.json();

  // add or remove ticker from watchlist
  if (action === "add") {
    await addToWatchlist(tickers);
  } else if (action === "remove") {
    await removeFromWatchlist(tickers);
  }

  // return success
  return new Response(JSON.stringify({ success: true }), {
    status: 200,
    headers: {
      "Content-Type": "application/json",
    },
  });
};
