import { parseQuery } from "../utils/query-parser";
import { Fn, HTTP_METHODS } from "./http";

class RouteNode {
  children: Map<string, RouteNode>;
  handler: Map<string, Fn[]>;
  middlewares: Fn[] = [];
  params: string[];

  constructor() {
    this.children = new Map();
    this.handler = new Map();
    this.params = [];
  }

  get(method: string) {
    return this.handler.get(method) || [];
  }
}

export class Route {
  private root: RouteNode;

  constructor() {
    this.root = new RouteNode();
  }

  middlewares() {
    return this.root.middlewares;
  }

  get(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.GET, ...handler);
  }

  post(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.POST, ...handler);
  }

  put(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.PUT, ...handler);
  }

  delete(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.DELETE, ...handler);
  }

  patch(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.PATCH, ...handler);
  }

  head(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.HEAD, ...handler);
  }

  options(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.OPTIONS, ...handler);
  }

  connect(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.CONNECT, ...handler);
  }

  trace(path: string, ...handler: Fn[]) {
    return this.addRoute(path, HTTP_METHODS.TRACE, ...handler);
  }

  use(...args: [path: string, router: Route] | Fn[]) {
    if (typeof args[0] === "string" && args[1] instanceof Route) {
      const [path, router] = args as [string, Route];
      this.mountRouter(path, router);
      return;
    }

    const middlewares = args as Fn[];
    this.root.middlewares.push(...middlewares);
  }

  findRoute(
    path: string = "/",
    method: string = HTTP_METHODS.GET
  ): {
    params: Record<string, string>;
    query: Record<string, string | string[]>;
    handler: Fn[];
  } {
    let params: Record<string, string> = {};
    let query: Record<string, string | string[]> = {};
    let handler = new Array<Fn>();

    if (path.length == 0) return { params, query, handler };
    if (path == "/") return { params, query, handler: this.root.get(method) };

    let endPathIdx = path.length;
    if (path.indexOf("?") !== -1) {
      endPathIdx = path.indexOf("?");
    }

    const words = path.substring(0, endPathIdx).split("/").filter(Boolean);
    let curNode = this.root;
    let segments: string[] = [];

    for (const word of words) {
      const segment = word.trim().toLowerCase();
      const children = curNode.children;

      if (children.has(segment)) {
        curNode = children.get(segment)!;
        handler = curNode.get(method);
      } else if (children.has(":")) {
        // dynamic
        curNode = children.get(":")!;
        handler = curNode.get(method);
        segments.push(segment);
      } else {
        return { params, query, handler };
      }
    }

    for (let i = 0; i < segments.length; i++) {
      params[curNode.params[i]] = segments[i];
    }

    query = parseQuery(path);
    return { params, query, handler };
  }

  printTree(node = this.root, indentation = 0) {
    const indent = "-".repeat(indentation);
    node.children.forEach((childNode, segment) => {
      console.log(`${indent}(${segment}) Dynamic: ${childNode.params}`);
      this.printTree(childNode, indentation + 1);
    });
  }

  private addRoute(path: string, method: string, ...handler: Fn[]) {
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

  private mountRouter(path: string, router: Route) {
    const segments = path.split("/").filter(Boolean);
    let current = this.root;

    for (const segment of segments) {
      const key = segment.startsWith(":") ? ":" : segment;
      if (!current.children.has(key)) {
        current.children.set(key, new RouteNode());
      }
      current = current.children.get(key)!;
      if (segment.startsWith(":")) {
        current.params.push(segment.slice(1));
      }
    }

    const stack: { from: RouteNode; to: RouteNode }[] = [
      { from: router.root, to: current },
    ];

    while (stack.length > 0) {
      const { from, to } = stack.pop()!;

      to.middlewares.push(...from.middlewares);

      from.handler.forEach((handlers, method) => {
        if (!to.handler.has(method)) to.handler.set(method, []);
        to.handler.get(method)!.push(...handlers);
      });

      from.children.forEach((childFrom, key) => {
        if (!to.children.has(key)) {
          to.children.set(key, new RouteNode());
        }
        const childTo = to.children.get(key)!;
        stack.push({ from: childFrom, to: childTo });
      });

      to.params = from.params.slice();
    }
  }
}
