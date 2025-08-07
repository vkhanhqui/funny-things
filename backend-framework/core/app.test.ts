import { describe, it, expect, beforeAll, afterAll } from "vitest";
import request from "supertest";
import { App } from "./app";
import { Route } from "./route";
import { Server } from "http";
import { HttpError } from "../internal/error";
import { logging } from "../internal/middleware";

describe("App", () => {
  let app: App;

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

    const userRouter = new Route();
    userRouter.get("/400", (req, res, next) => {
      next(HttpError.BadRequest("Bad request"));
    });
    router.use("/users", userRouter);

    app = new App(router);
    await app.listen();
  });

  it("GET /200 should return 200 OK", async () => {
    await request(app.getServer())
      .get("/200")
      .expect(200)
      .expect("Content-Type", /text\/plain/)
      .expect("Hello from the root endpoint");
  });

  it("GET /404 should return 404 Not Found", async () => {
    await request(app.getServer())
      .get("/404")
      .expect(404)
      .expect("Content-Type", /json/)
      .expect({ error: "Not Found" });
  });

  it("GET /400 should return 400 Bad Request", async () => {
    await request(app.getServer())
      .get("/400")
      .expect(400)
      .expect("Content-Type", /json/)
      .expect({ error: "Bad request" });
  });

  it("GET /500 should return 500 Internal Server Error", async () => {
    await request(app.getServer())
      .get("/500")
      .expect(500)
      .expect("Content-Type", /json/)
      .expect({ error: "Unexpected error" });
  });

  it("GET /500/async should return 500 Internal Server Error", async () => {
    await request(app.getServer())
      .get("/500/async")
      .expect(500)
      .expect("Content-Type", /json/)
      .expect({ error: "Unexpected error" });
  });

  it("GET /users/400 should return 400 Bad Request", async () => {
    await request(app.getServer())
      .get("/users/400")
      .expect(400)
      .expect("Content-Type", /json/)
      .expect({ error: "Bad request" });
  });

  it("GET /users/404 should return 404 Not Found", async () => {
    await request(app.getServer())
      .get("/users/404")
      .expect(404)
      .expect("Content-Type", /json/)
      .expect({ error: "Not Found" });
  });
});
