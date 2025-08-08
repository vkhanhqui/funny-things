import { Route, App } from "./core";
import { logging, jsonParser } from "./internal/middleware";
import { HttpError } from "./internal/error";

const router = new Route();
router.get("/200", (req, res, next) => {
  res.writeHead(200, { "Content-Type": "text/plain" });
  res.end("Hello from the root endpoint");
  next();
});

router.get("/400", (req, res, next) => {
  next(HttpError.BadRequest("Bad request"));
});

router.get("/500", (req, res, next) => {
  throw new Error("Unexpected error");
});

router.get("/500/async", async (req, res, next) => {
  throw new Error("Unexpected error");
});

router.post("/json", (req, res, next) => {
  res.writeHead(200, { "Content-Type": "application/json" });
  res.end(req.body);
});

const userRouter = new Route();
userRouter.get("/400", (req, res, next) => {
  next(HttpError.BadRequest("Bad request"));
});
// router.use(logging);
router.use("/users", userRouter);

const app = new App(router);
app.use(logging);
app.use(jsonParser());
app.listen(3000);
