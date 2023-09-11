// YOUR_BASE_DIRECTORY/netlify/functions/test-scheduled-function.ts

import type { Handler, HandlerEvent, HandlerContext } from "@netlify/functions";
import { schedule } from "@netlify/functions";

const myHandler: Handler = async (event: HandlerEvent, context: HandlerContext) => {
  console.log("[ INFO ] Inside scheduled-minutely.ts");
  console.log("[ INFO ] Received event:", event);

  return {
    statusCode: 200,
  };
};

// every day at 4:20pm
const handler = schedule("20 16 * * *", myHandler);

export { handler };