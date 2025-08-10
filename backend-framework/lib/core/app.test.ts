import { describe, it, beforeAll } from "vitest";
import request from "supertest";
import { App } from "./app";
import { Route } from "./route";
import { Server } from "http";
import { HttpError } from "../error";
import { jsonParser, logging, compression } from "../middleware";
import { Req, Res } from "./http";

describe("App", () => {
  let server: Server<typeof Req, typeof Res>;

  beforeAll(async () => {
    const router = new Route();
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

    router.get("/params/:id/:name", (req, res, next) => {
      res.json(200, req.params);
    });

    router.get("/query", (req, res, next) => {
      res.json(200, req.query);
    });

    const userRouter = new Route();
    userRouter.get("/400", (req, res, next) => {
      next(HttpError.BadRequest("Bad request"));
    });
    router.use("/users", userRouter);

    const app = new App(router);
    app.use(logging);
    app.use(jsonParser());
    // app.use(compression(0));

    const onListening = new Promise<Server<typeof Req, typeof Res>>(
      (resolve) => {
        const onListeningResolve = () => resolve(app.getServer());
        app.listen(0, onListeningResolve);
      }
    );
    server = await onListening;
  });

  it("GET /200 should return 200 OK", async () => {
    await request(server)
      .get("/200")
      .expect(200)
      .expect("Content-Type", /text\/html/)
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

  it("GET /params should return the same params", async () => {
    await request(server)
      .get("/params/123/quivo")
      .set("Content-Type", "application/json")
      .expect(200)
      .expect("Content-Type", /json/)
      .expect({ id: "123", name: "quivo" });
  });

  it("GET /query should return the same query", async () => {
    await request(server)
      .get("/query?name=John&age=25")
      .set("Content-Type", "application/json")
      .expect(200)
      .expect("Content-Type", /json/)
      .expect({ name: "John", age: "25" });
  });

  it("Not modified response with Etag should return 304", async () => {
    await request(server)
      .get("/query?name=John&age=25")
      .set("Content-Type", "application/json")
      .expect(200)
      .expect("Content-Type", /json/)
      .expect({ name: "John", age: "25" });

    await request(server)
      .get("/query?name=John&age=25")
      .set("Content-Type", "application/json")
      .set("if-none-match", `W/"OLcYrp4Qxdu+YFfyhT96NbtEx28="`)
      .expect(304)
      .expect({});
  });
});
