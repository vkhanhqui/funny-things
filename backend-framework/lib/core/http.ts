import { IncomingMessage, ServerResponse } from "http";
import { createHash } from "crypto";

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

export class Res extends ServerResponse {
  send(code: number, body: any, replyWith304 = true) {
    var chunk = body;
    switch (typeof chunk) {
      case "string":
        if (!this.getContentType()) {
          this.setContentType("text/html");
        }
        break;
      case "boolean":
      case "number":
      case "object":
        if (chunk === null) {
          chunk = "";
          break;
        }
        if (ArrayBuffer.isView(chunk)) {
          if (!this.getContentType()) {
            this.setContentType("application/octet-stream");
          }
          break;
        }
        return this.json(code, chunk);
      default:
        throw Error(`Not support body type: ${typeof body}`);
    }

    if (replyWith304) {
      this.handleETag(body);
    }
    if (this.writableEnded) {
      return;
    }

    this.statusCode = code;
    this.end(chunk);
  }

  json(code: number, body: any) {
    if (!this.getContentType()) {
      this.setContentType("application/json");
    }

    return this.send(code, JSON.stringify(body));
  }

  private getContentType() {
    return this.getHeader("Content-Type");
  }

  private setContentType(value: string) {
    return this.setHeader("Content-Type", value);
  }

  private handleETag(body: any) {
    let etag = this.getHeader("etag");
    if (!etag) {
      etag = this.generateWeakETag(body as string | Buffer);
      this.setHeader("etag", etag);
    }

    var ifNoneMatch = this.req.headers["if-none-match"];
    if (
      ifNoneMatch === etag ||
      ifNoneMatch === `W/"${etag}"` ||
      `W/"${ifNoneMatch}"` === etag
    ) {
      this.statusCode = 304;
      this.removeHeader("Content-Type");
      this.removeHeader("Content-Length");
      this.removeHeader("Transfer-Encoding");
      this.end("");
    }
  }

  private generateWeakETag(body: string | Buffer) {
    return `W/"${createHash("sha1").update(body).digest("base64")}"`;
  }
}

export type Fn = (
  req: Req,
  res: Res,
  next: (err?: any) => void
) => void | Promise<void>;
