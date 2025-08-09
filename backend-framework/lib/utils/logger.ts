enum LogLevel {
  INFO = "INFO",
  WARN = "WARN",
  ERROR = "ERROR",
  DEBUG = "DEBUG",
}

class Logger {
  private format(level: LogLevel, message: string, meta?: Record<string, any>) {
    const timestamp = new Date().toISOString();
    const logObj = {
      level,
      timestamp,
      message,
      ...(meta || {}),
    };

    return JSON.stringify(logObj);
  }

  info(message: string, meta?: Record<string, any>) {
    console.log(this.format(LogLevel.INFO, message, meta));
  }

  warn(message: string, meta?: Record<string, any>) {
    console.warn(this.format(LogLevel.WARN, message, meta));
  }

  error(message: string, meta?: Record<string, any>) {
    console.error(this.format(LogLevel.ERROR, message, meta));
  }

  debug(message: string, meta?: Record<string, any>) {
    console.debug(this.format(LogLevel.DEBUG, message, meta));
  }
}

export const logger = new Logger();
