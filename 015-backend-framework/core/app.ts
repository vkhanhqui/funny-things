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
      const next = () => {
        if (i >= handlers.length) {
          return;
        }
        const handler = handlers[i++];
        handler(req, res, next);
      };

      next();
    }).listen(port, () => {
      console.log(`Server running at http://localhost:${port}`);
    });
  }
}
