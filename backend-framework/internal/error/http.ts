export class HttpError extends Error {
  statusCode: number;

  constructor(message: string, statusCode = 500) {
    super(message);
    this.statusCode = statusCode;
    Object.setPrototypeOf(this, HttpError.prototype);
  }

  static BadRequest(message = "Bad Request") {
    return new HttpError(message, 400);
  }

  static Unauthorized(message = "Unauthorized") {
    return new HttpError(message, 401);
  }

  static Forbidden(message = "Forbidden") {
    return new HttpError(message, 403);
  }

  static NotFound(message = "Not Found") {
    return new HttpError(message, 404);
  }
}
