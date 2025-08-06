import { createServer, IncomingMessage, ServerResponse } from "http";
import { Fn, Route } from "./route";

export class App {
  private router: Route;
  private middlewares: Fn[] = [];

  constructor(router: Route) {
    this.router = router;
  }

  use(middleware: Fn) {
    this.middlewares.push(middleware);
  }

  listen(port: number) {
    createServer((req: IncomingMessage, res: ServerResponse) => {
      const route = this.router.findRoute(req.url, req.method);
      const handlers = [...this.middlewares, ...route.handler];

      let i = 0;
      const next = (err?: any) => {
        if (err) {
          this.handleError(req, res, err)
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
    }).listen(port, () => {
      console.log(`Server running at port ${port}`);
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
