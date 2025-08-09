import { describe, test, expect } from "vitest";
import { parseQuery } from "./query-parser";

describe("parseQuery", () => {
  test("returns empty object when no query string", () => {
    expect(parseQuery("/")).toEqual({});
  });

  test("parses single query parameter", () => {
    expect(parseQuery("/?name=John")).toEqual({ name: "John" });
  });

  test("parses multiple different query parameters", () => {
    expect(parseQuery("/?name=John&age=25")).toEqual({
      name: "John",
      age: "25",
    });
  });

  test("parses array", () => {
    expect(parseQuery("/?tag=[js,ts,go]")).toEqual({
      tag: ["js", "ts", "go"],
    });
  });

  test("parses duplicate keys into array", () => {
    expect(parseQuery("/?tag=js&tag=ts&tag=go")).toEqual({
      tag: ["js", "ts", "go"],
    });
  });

  test("handles empty value", () => {
    expect(parseQuery("/?name=")).toEqual({ name: "" });
  });

  test("handles missing '=' in param", () => {
    expect(parseQuery("/?name")).toEqual({ name: "" });
  });

  test("handles path without leading slash", () => {
    expect(parseQuery("?name=John")).toEqual({ name: "John" });
  });

  test("handles encoded characters", () => {
    expect(parseQuery("/?name=John%20Doe")).toEqual({
      name: "John Doe",
    });
  });

  test("handles multiple '&' without values", () => {
    expect(parseQuery("/?name=John&&age=25")).toEqual({
      name: "John",
      age: "25",
    });
  });
});
