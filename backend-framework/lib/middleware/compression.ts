import zlib from "zlib";
import { Req, Res } from "../core/http";

const DEFAULT_THRESHOLD = 1024; // 1KB

export function compression(threshold: number = DEFAULT_THRESHOLD) {
  return (req: Req, res: Res, next: (err?: any) => void) => {
    const encoding = pickEncoding(req);
    if (!encoding) {
      next();
      return;
    }

    const _writeHead = res.writeHead.bind(res);
    const _write = res.write.bind(res);
    const _end = res.end.bind(res);
    let firstChecked = false;

    let compressor: zlib.BrotliCompress | zlib.Gzip | zlib.Deflate;
    function startCompression() {
      compressor = createCompressor(encoding);
      compressor.on("data", (chunk) => _write(chunk));
      compressor.on("end", () => _end());
    }

    function writeHead(chunk: any, encodingArg?: any) {
      if (firstChecked) {
        return;
      }
      firstChecked = true;
      if (chunkLength(chunk, encodingArg) >= threshold) {
        cpHeaders = { ...cpHeaders, "Content-Encoding": encoding };
        startCompression();
      }
      _writeHead(cpStatusCode, cpHeaders);
    }

    let cpStatusCode: number;
    let cpHeaders: any;

    res.writeHead = (
      statusCode: number,
      statusMessage?: any,
      headers?: any
    ) => {
      cpStatusCode = statusCode;
      cpHeaders = { ...headers, ...statusMessage };
      return res;
    };

    res.write = (chunk: any, encodingArg?: any, cb?: any) => {
      writeHead(chunk, encodingArg);

      if (compressor) {
        compressor.write(chunk, encodingArg, cb);
      } else {
        _write(chunk, encodingArg, cb);
      }
      return true;
    };

    res.end = (chunk?: any, encodingArg?: any, cb?: any) => {
      writeHead(chunk, encodingArg);

      if (compressor) {
        if (chunk) compressor.end(chunk, encodingArg, cb);
        else compressor.end();
      } else {
        if (chunk) _end(chunk, encodingArg, cb);
        else _end();
      }
      return res;
    };

    next();
  };
}

function pickEncoding(req: Req) {
  const acceptEncoding = req.headers["accept-encoding"] || "";
  if (!acceptEncoding) return "";
  if (acceptEncoding.includes("br")) return "br";
  if (acceptEncoding.includes("gzip")) return "gzip";
  if (acceptEncoding.includes("deflate")) return "deflate";
  return "";
}

function createCompressor(encoding: string) {
  switch (encoding) {
    case "br":
      return zlib.createBrotliCompress();
    case "gzip":
      return zlib.createGzip();
    case "deflate":
      return zlib.createDeflate();
    default:
      return zlib.createBrotliCompress();
  }
}

function chunkLength(chunk: any, bufferEncoding: BufferEncoding) {
  if (!chunk) return 0;

  return Buffer.isBuffer(chunk)
    ? chunk.length
    : Buffer.byteLength(chunk, bufferEncoding);
}
