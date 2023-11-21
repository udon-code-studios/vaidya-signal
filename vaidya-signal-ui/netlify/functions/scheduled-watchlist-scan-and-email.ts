import type { Handler, HandlerEvent, HandlerContext } from "@netlify/functions";
import { schedule } from "@netlify/functions";

const vaidyaServiceUrl = "https://vaidya-service.udon.studio";

const myHandler: Handler = async (event: HandlerEvent, context: HandlerContext) => {
  console.log("[ INFO ] Inside scheduled-watchlist-scan-and-email.ts");
  console.log("[ INFO ] Time:", new Date());

  await fetch(`${vaidyaServiceUrl}/api/v1/update-watchlist-email-today-triggers`, {
    method: "POST",
  });

  return {
    statusCode: 200,
  };
};

// every day at 4:20pm EST
const handler = schedule("20 21 * * *", myHandler);

export { handler };
