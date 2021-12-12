import wretch from "wretch";
import { dedupe, retry, throttlingCache } from "wretch-middlewares";

export const wretchInstance = wretch()
    .url("api/v0/")
    .middlewares([dedupe(), retry(), throttlingCache()]);
