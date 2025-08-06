import { Route, App } from './core';
import { logging } from './internal/middleware';

const router = new Route();
router.get("/", (req, res, next) => {
  res.writeHead(200, { "Content-Type": "text/plain" });
  res.end("Hello from the root endpoint");
  next();
});

const app = new App(router);
app.use(logging);
app.listen(3000);