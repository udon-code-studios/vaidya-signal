import type { Handler, HandlerEvent, HandlerContext } from "@netlify/functions";
import { schedule } from "@netlify/functions";

const vaidyaServiceUrl = import.meta.env.VAIDYA_SERVICE_URL;

const myHandler: Handler = async (event: HandlerEvent, context: HandlerContext) => {
  console.log("[ INFO ] Inside scheduled-watchlist-scan-and-email.ts");
  console.log("[ INFO ] Time:", new Date());

  fetch(`${vaidyaServiceUrl}/api/v1/update-watchlist-email-today-triggers`, {
    method: "POST",
  });

  return {
    statusCode: 200,
  };
};

// every day at 4:20pm
const handler = schedule("20 16 * * *", myHandler);

export { handler };
