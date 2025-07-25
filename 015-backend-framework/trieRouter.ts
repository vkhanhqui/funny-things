import { createServer, IncomingMessage, ServerResponse } from "http";

const HTTP_METHODS = {
  GET: "GET",
  POST: "POST",
  PUT: "PUT",
  DELETE: "DELETE",
  PATCH: "PATCH",
  HEAD: "HEAD",
  OPTIONS: "OPTIONS",
  CONNECT: "CONNECT",
  TRACE: "TRACE",
};

class RouteNode {
  children: Map<string, RouteNode>;
  handler: Map<string, Function>;
  params: string[];

  constructor() {
    this.children = new Map();
    this.handler = new Map();
    this.params = [];
  }
}

class Route {
  root: RouteNode;

  constructor() {
    this.root = new RouteNode();
  }

  get(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.GET, handler);
  }

  post(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.POST, handler);
  }

  put(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.PUT, handler);
  }

  delete(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.DELETE, handler);
  }

  patch(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.PATCH, handler);
  }

  head(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.HEAD, handler);
  }

  options(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.OPTIONS, handler);
  }

  connect(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.CONNECT, handler);
  }

  trace(path: string, handler: Function) {
    return this.addRoute(path, HTTP_METHODS.TRACE, handler);
  }

  findRoute(
    path: string,
    method: string
  ): { params: Map<string, string>; handler: Function | undefined } {
    const params = new Map();
    let handler: Function | undefined = undefined;
    if (path.length == 0) return { params, handler };
    if (path == "/") return { params, handler: this.root.handler.get(method) };

    const words = path.split("/").filter(Boolean);
    let curNode = this.root;
    let segments: string[] = [];

    for (const word of words) {
      const segment = word.trim().toLowerCase();
      const children = curNode.children;

      if (children.has(segment)) {
        curNode = children.get(segment)!;
        handler = curNode.handler.get(method);
      } else if (children.has(":")) {
        // dynamic
        curNode = children.get(":")!;
        handler = curNode.handler.get(method);
        segments.push(segment);
      } else {
        return { params, handler };
      }
    }

    for (let i = 0; i < segments.length; i++) {
      params.set(curNode.params[i], segments[i]);
    }

    return { params, handler };
  }

  printTree(node = this.root, indentation = 0) {
    const indent = "-".repeat(indentation);
    node.children.forEach((childNode, segment) => {
      console.log(`${indent}(${segment}) Dynamic: ${childNode.params}`);
      this.printTree(childNode, indentation + 1);
    });
  }

  private addRoute(path: string, method: string, handler: Function) {
    this.verifyParams(path, method);

    let cur = this.root;
    let dynamicParams: string[] = [];
    const words = path.split("/").filter(Boolean);

    for (let i = 0; i < words.length; i++) {
      const word = words[i];
      const isDynamic = word[0] == ":";
      const segment = isDynamic ? ":" : word.trim().toLowerCase();

      if (isDynamic) {
        dynamicParams.push(word.substring(1));
      }

      if (!cur.children.has(segment)) {
        cur.children.set(segment, new RouteNode());
      }
      cur = cur.children.get(segment)!;
    }

    cur.handler.set(method, handler);
    cur.params = dynamicParams;
  }

  private verifyParams(path: string, method: string) {
    if (path.length == 0 || path[0] !== "/") {
      throw new Error("Malformed path");
    }

    if (!HTTP_METHODS[method]) {
      throw new Error("Invalid HTTP method");
    }
  }
}

function run(router: Route, port: number) {
  createServer((req: IncomingMessage, res: ServerResponse) => {
    const route = router.findRoute(req.url || "/", req.method || "GET");
    if (route.handler) {
      route.handler(req, res);
    } else {
      res.writeHead(404, { "Content-Type": "text/plain" });
      res.end("Not Found");
    }
  }).listen(port, () => {
    console.log(`Server listening at http://localhost:${port}`);
  });
}

const router = new Route();
router.get("/", (req: IncomingMessage, res: ServerResponse) => {
  res.writeHead(200, { "Content-Type": "text/plain" });
  res.end("Hello from the root endpoint");
});
router.printTree();

run(router, 3000);
