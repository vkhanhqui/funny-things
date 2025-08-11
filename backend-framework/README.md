# Build a backend web framework from scratch

## Quick Start

```ts
import { Route, App } from "./lib";

const router = new Route();
router.get("/", (req, res) => {
  res.send(200, "Hello World!");
});

const app = new App(router);
app.listen(3000, () => {
  console.log("Server running at http://localhost:3000");
});
```

## Basic Usage
```ts
import { Route, App } from "./lib";
import { jsonParser, logging } from "./lib/middleware";
import { HttpError } from "./lib/error";

// Create a router
const router = new Route();

// Define routes
router.get("/", (req, res) => {
  res.send(200, "Hello World!");
});

router.get("/users/:id", (req, res) => {
  res.json(200, { userId: req.params!["id"] });
});

router.post("/echo", (req, res) => {
  res.json(200, req.body);
});

// Throwing an HTTP error
router.get("/bad-request", (req, res, next) => {
  next(HttpError.BadRequest("Invalid input"));
});

// Create the app
const app = new App(router);

// Register middlewares
app.use(logging);
app.use(jsonParser());

// Start listening
app.listen(3000, () => {
  console.log("Server is running on http://localhost:3000");
});
```

## Nested Routers
```ts
const apiRouter = new Route();
const userRouter = new Route();

userRouter.get("/", (req, res) => {
  res.json(200, [{ id: 1, name: "John Doe" }]);
});

userRouter.get("/:id", (req, res) => {
  res.json(200, { id: req.params!["id"], name: "John Doe" });
});

apiRouter.use("/users", userRouter);

const app = new App(apiRouter);
app.listen(3000);
```

## Middleware
```ts
import { logging, jsonParser, compression } from "./lib/middleware";

app.use(logging);
app.use(jsonParser());
// app.use(compression());
```

## Error Handling
```ts
import { HttpError } from "./lib/error";

router.get("/error", () => {
  throw HttpError.InternalServerError("Something went wrong");
});
```

## Features
- HTTP server
- Dynamic routes
- Middleware
- Logging
- Error handling
- Route group
- Body parsing (json for now)
- Path variable + Query params
- HTTP compression
- ETags for HTTP responses

## In progress
- Fix HTTP compression bug

## TODO
- Caching
- Accepts in requests
- Request/response validation

## References
- [Learn nodejs the hard way](https://github.com/ishtms/learn-nodejs-hard-way)
- [expressjs.com](https://expressjs.com/)
- [fastify](https://github.com/fastify/fastify/blob/main/docs/Guides/Ecosystem.md#core)