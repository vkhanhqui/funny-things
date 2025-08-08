import { IncomingMessage, ServerResponse } from "http";

export class Req extends IncomingMessage {
  body?: any;
}

export class Res extends ServerResponse {}
