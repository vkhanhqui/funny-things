export function parseQuery(path: string): Record<string, string | string[]> {
  let query: Record<string, string | string[]> = {};
  const idx = path.indexOf("?");
  if (idx === -1) {
    return query;
  }

  const queryPath = path.substring(idx + 1).split("&");
  for (const pair of queryPath) {
    if (!pair) continue;

    const [rawK, rawV] = pair.split("=");
    const k = decodeURIComponent(rawK || "");
    const v = decodeURIComponent(rawV || "");

    if (v.startsWith("[") && v.endsWith("]")) {
      const items = v
        .slice(1, -1)
        .split(",")
        .map((v) => v.trim());
      query[k] = items;
      continue;
    }

    const existing = query[k];
    if (existing === undefined) {
      query[k] = v;
    } else if (Array.isArray(existing)) {
      existing.push(v);
    } else {
      query[k] = [existing, v];
    }
  }
  return query;
}
