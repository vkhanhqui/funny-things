import { Req, Res } from "../core/http";
import { logger } from "../utils";
import { randomUUID } from "crypto";

export const logging = (
  req: Req,
  res: Res,
  next: (err?: any) => void
) => {
  const start = process.hrtime.bigint();
  const requestId = randomUUID();

  const ip = req.socket.remoteAddress;
  const userAgent = req.headers["user-agent"];

  logger.info("Request", {
    "type": "request",
    "request_id": requestId,
    "user_agent": userAgent,
    method: req.method,
    url: req.url,
    ip,
  });

  const originalEnd = res.end;
  res.end = function (chunk?: any, encoding?: any, cb?: any) {
    const durationNs = process.hrtime.bigint() - start;
    const durationMs = Number(durationNs) / 1_000_000;

    logger.info("Response", {
      "type": "response",
      "request_id": requestId,
      "status_code": res.statusCode,
      "response_time": durationMs.toFixed(2),
      "user_agent": userAgent,
      method: req.method,
      url: req.url,
      ip,
    });

    return originalEnd.call(this, chunk, encoding, cb);
  };

  next();
};
