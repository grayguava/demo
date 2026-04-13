#!/usr/bin/env node
"use strict";

const { spawnSync } = require("child_process");
const path          = require("path");

const MIN_MAJOR = 3;
const MIN_MINOR = 8;

const SCRIPT = path.join(__dirname, "fsync.py");
const ARGS   = process.argv.slice(2);

function findPython() {
  const candidates = ["py", "python", "python3"];

  for (const cmd of candidates) {
    const result = spawnSync(cmd, ["--version"], { encoding: "utf8" });
    if (result.status === 0) {
      const output = (result.stdout || result.stderr || "").trim();
      const match  = output.match(/Python (\d+)\.(\d+)/);
      if (match) {
        const major = parseInt(match[1], 10);
        const minor = parseInt(match[2], 10);
        if (major === MIN_MAJOR && minor >= MIN_MINOR) return { cmd, version: output };
        if (major > MIN_MAJOR) return { cmd, version: output };
        fail(
          "Python " + MIN_MAJOR + "." + MIN_MINOR + "+ is required.\n" +
          "           Found: " + output + "\n" +
          "           Install from https://python.org"
        );
      }
    }
  }
  return null;
}

function fail(msg) {
  console.error("\x1b[31m\x1b[1m ERR \x1b[0m " + msg);
  process.exit(1);
}

const python = findPython();

if (!python) {
  fail(
    "Python is required to run fsync.\n" +
    "           Install from https://python.org\n" +
    "           Windows: winget install Python.Python.3"
  );
}

const result = spawnSync(python.cmd, [SCRIPT, ...ARGS], {
  stdio: "inherit",
  env:   process.env,
});

process.exit(result.status ?? 1);