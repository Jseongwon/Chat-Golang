const { createLogger, format, transports } = require("winston");
const winstonDaily = require("winston-daily-rotate-file");

const logDir = "logs";

const logFormat = format.printf((info) => {
  return `${info.timestamp} ${info.level}: ${info.message}`;
});

const SocketLogger = createLogger({
  format: format.combine(
    format.timestamp({
      format: "YYYY-MM-DD HH:mm:ss",
    }),
    logFormat
  ),
  transports: [
    new winstonDaily({
      filename: `%DATE%.log`,
      datePattern: "YYYY-MM-DD",
      dirname: logDir + "/socket-info",
      level: "info",
    }),

    new winstonDaily({
      filename: `%DATE%.log`,
      datePattern: "YYYY-MM-DD",
      dirname: logDir + "/socket-error",
      level: "error",
    }),
  ],
});

SocketLogger.add(
  new transports.Console({
    format: format.combine(format.colorize(), format.simple()),
  })
);

module.exports = { SocketLogger };
