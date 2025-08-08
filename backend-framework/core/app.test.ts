import { describe, it, beforeAll } from "vitest";
import request from "supertest";
import { App } from "./app";
import { Route } from "./route";
import { Server } from "http";
import { HttpError } from "../internal/error";
import { jsonParser, logging } from "../internal/middleware";

describe("App", () => {
  let server: Server;

  beforeAll(async () => {
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
      next();
    });

    const userRouter = new Route();
    userRouter.get("/400", (req, res, next) => {
      next(HttpError.BadRequest("Bad request"));
    });
    router.use("/users", userRouter);

    const app = new App(router);
    app.use(logging);
    app.use(jsonParser());

    const onListening = new Promise<Server>((resolve) => {
      const onListeningResolve = () => resolve(app.getServer());
      app.listen(0, onListeningResolve);
    });
    server = await onListening;
  });

  it("GET /200 should return 200 OK", async () => {
    await request(server)
      .get("/200")
      .expect(200)
      .expect("Content-Type", /text\/plain/)
      .expect("Hello from the root endpoint");
  });

  it("GET /404 should return 404 Not Found", async () => {
    await request(server)
      .get("/404")
      .expect(404)
      .expect("Content-Type", /json/)
      .expect({ error: "Not Found" });
  });

  it("GET /400 should return 400 Bad Request", async () => {
    await request(server)
      .get("/400")
      .expect(400)
      .expect("Content-Type", /json/)
      .expect({ error: "Bad request" });
  });

  it("GET /500 should return 500 Internal Server Error", async () => {
    await request(server)
      .get("/500")
      .expect(500)
      .expect("Content-Type", /json/)
      .expect({ error: "Unexpected error" });
  });

  it("GET /500/async should return 500 Internal Server Error", async () => {
    await request(server)
      .get("/500/async")
      .expect(500)
      .expect("Content-Type", /json/)
      .expect({ error: "Unexpected error" });
  });

  it("GET /users/400 should return 400 Bad Request", async () => {
    await request(server)
      .get("/users/400")
      .expect(400)
      .expect("Content-Type", /json/)
      .expect({ error: "Bad request" });
  });

  it("GET /users/404 should return 404 Not Found", async () => {
    await request(server)
      .get("/users/404")
      .expect(404)
      .expect("Content-Type", /json/)
      .expect({ error: "Not Found" });
  });

  it("POST /json should return the same JSON body", async () => {
    const payload = { message: "Hello from test" };
    await request(server)
      .post("/json")
      .send(payload)
      .set("Content-Type", "application/json")
      .expect(200)
      .expect("Content-Type", /json/)
      .expect(payload);
  });
});
