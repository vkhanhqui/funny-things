import { HttpError } from "../error";
import { Req, Res, Fn } from "../core/http";

export function jsonParser() {
  return async (
    req: Req,
    res: Res,
    next: (err?: any) => void
  ) => {
    if (
      req.method === "GET" ||
      req.method === "HEAD" ||
      req.headers["content-type"] !== "application/json"
    ) {
      return next();
    }

    const chunks: any[] = [];

    req.on("data", (chunk) => {
      chunks.push(chunk);
    });

    req.on("end", () => {
      try {
        req.body = Buffer.concat(chunks).toString();
        next();
      } catch (err) {
        next(HttpError.BadRequest("Invalid JSON"));
      }
    });

    req.on("error", () => {
      next(HttpError.BadRequest("Error parsing request body"));
    });
  };
}
