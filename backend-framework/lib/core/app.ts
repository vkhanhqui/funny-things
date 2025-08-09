import { createServer, IncomingMessage, Server, ServerResponse } from "http";
import { Route } from "./route";
import { HttpError } from "../error";
import { Req, Res, Fn } from "./http";

export class App {
  private server: Server;
  private router: Route;
  private middlewares: Fn[] = [];

  constructor(router: Route) {
    this.router = router;
  }

  use(middleware: Fn) {
    this.middlewares.push(middleware);
  }

  listen(port: number = 0, onListening?: () => void) {
    this.server = this.createServer();
    this.server.listen(port, () => {
      console.log(`Server running at port ${port}`);
      onListening?.();
    });
  }

  getServer() {
    return this.server;
  }

  private createServer() {
    return createServer((req: Req, res: Res) => {
      const route = this.router.findRoute(req.url, req.method);
      const handlers = [
        ...this.middlewares,
        ...this.router.middlewares(),
        ...route.handler,
      ];

      req.params = { ...route.params };

      let i = 0;
      const next = (err?: any) => {
        if (err) {
          this.handleError(req, res, err);
          return;
        }

        if (route.handler.length == 0) {
          next(HttpError.NotFound("Not Found"));
          return;
        }

        if (i >= handlers.length) {
          return;
        }

        const handler = handlers[i++];
        try {
          const v = handler(req, res, next);
          if (v && typeof v.then === "function") {
            Promise.resolve(v).catch(next);
          }
        } catch (error) {
          next(error);
        }
      };

      next();
    });
  }

  private handleError(req: IncomingMessage, res: ServerResponse, err: any) {
    if (res.writableEnded) return;

    const statusCode = err.statusCode || 500;
    const message = err.message || "Internal Server Error";

    res.writeHead(statusCode, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ error: message }));
  }
}
