enum LogLevel {
  INFO = "INFO",
  WARN = "WARN",
  ERROR = "ERROR",
  DEBUG = "DEBUG",
}

class Logger {
  private format(level: LogLevel, message: any, meta?: Record<string, any>) {
    const timestamp = new Date().toISOString();
    const logObj = {
      level,
      timestamp,
      message,
      ...(meta || {}),
    };

    return JSON.stringify(logObj);
  }

  info(message: any, meta?: Record<string, any>) {
    console.log(this.format(LogLevel.INFO, message, meta));
  }

  warn(message: any, meta?: Record<string, any>) {
    const yellow = "\x1b[33m";
    const reset = "\x1b[0m";
    console.warn(
      `${yellow}${this.format(LogLevel.WARN, message, meta)}${reset}`
    );
  }

  error(message: any, meta?: Record<string, any>) {
    const red = "\x1b[31m";
    const reset = "\x1b[0m";
    console.error(
      `${red}${this.format(LogLevel.ERROR, message, meta)}${reset}`
    );
  }

  debug(message: any, meta?: Record<string, any>) {
    console.debug(this.format(LogLevel.DEBUG, message, meta));
  }
}

export const logger = new Logger();
