import { IncomingMessage, ServerResponse } from "http";

export const HTTP_METHODS = {
  GET: "GET",
  POST: "POST",
  PUT: "PUT",
  DELETE: "DELETE",
  PATCH: "PATCH",
  HEAD: "HEAD",
  OPTIONS: "OPTIONS",
  CONNECT: "CONNECT",
  TRACE: "TRACE",
};

export class Req extends IncomingMessage {
  body?: any = undefined;
  params?: Record<string, string> = {};
  query?: Record<string, string | string[]> = {};
}

export class Res extends ServerResponse {}

export type Fn = (
  req: Req,
  res: Res,
  next: (err?: any) => void
) => void | Promise<void>;
