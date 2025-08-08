import { IncomingMessage, ServerResponse } from "http";

export class Req extends IncomingMessage {
  body?: any = undefined;
  params?: Record<string, string> = {};
  query?: Record<string, string | string[]> = {};
}

export class Res extends ServerResponse {}
