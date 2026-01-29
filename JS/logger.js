const fs = require("fs");
const path = require("path");

function nowIso() {
  return new Date().toISOString();
}

function createLogger({ logsDir = "logs", filename }) {
  if (!fs.existsSync(logsDir)) fs.mkdirSync(logsDir, { recursive: true });

  const file =
    filename ||
    `game-${new Date()
      .toISOString()
      .replace(/[:.]/g, "-")
      .replace("T", "_")
      .slice(0, 19)}.jsonl`;

  const filepath = path.join(logsDir, file);

  function log(type, payload = {}) {
    const line = JSON.stringify({ ts: nowIso(), type, ...payload });
    fs.appendFileSync(filepath, line + "\n", "utf8");
  }

  return { filepath, log };
}

module.exports = { createLogger };
