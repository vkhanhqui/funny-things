import { Route, App } from "./lib";
import { logging, jsonParser, compression } from "./lib/middleware";
import { HttpError } from "./lib/error";

const router = new Route();

router.get("/params/:id/:name", (req, res, next) => {
  res.json(200, req.params);
});

router.get("/200", (req, res, next) => {
  res.send(200, "Hello from the root endpoint");
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
  res.json(200, req.body);
});

const userRouter = new Route();
userRouter.get("/400", (req, res, next) => {
  next(HttpError.BadRequest("Bad request"));
});
router.use("/users", userRouter);

const app = new App(router);
app.use(logging);
app.use(jsonParser());
// app.use(compression());
app.listen(3000);
