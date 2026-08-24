"use strict";

const fs = require("node:fs");
const path = require("node:path");
const { webcrypto } = require("node:crypto");

if (!globalThis.crypto) {
  globalThis.crypto = webcrypto;
}

require("./wasm_exec.js");

let invokeGo;

function initialize() {
  if (invokeGo) {
    return invokeGo;
  }

  const go = new globalThis.Go();
  const wasmPath = path.join(__dirname, "data-localizer.wasm");
  const module = new WebAssembly.Module(fs.readFileSync(wasmPath));
  const instance = new WebAssembly.Instance(module, go.importObject);
  void go.run(instance);

  if (typeof globalThis.__dataLocalizerInvoke !== "function") {
    throw new Error("the data-localizer Go WebAssembly runtime did not initialize");
  }
  invokeGo = globalThis.__dataLocalizerInvoke;
  return invokeGo;
}

function callGo(request) {
  let serialized;
  try {
    serialized = JSON.stringify(request);
  } catch (error) {
    throw new TypeError(`data must be JSON-compatible: ${error.message}`);
  }
  if (serialized === undefined) {
    throw new TypeError("data must be JSON-compatible");
  }

  const result = JSON.parse(initialize()(serialized));
  if (!result.ok) {
    throw hydrateError(result.error);
  }
  return result.value;
}

function localize(data, languageHeader = "", fallbackLanguage = "en") {
  return callGo({ data, languageHeader, fallbackLanguage });
}

function createLocalizer(options = {}) {
  const configuredOptions = {
    supportedLanguages: options.supportedLanguages,
    fallbackLanguage: options.fallbackLanguage,
    fallbackLanguages: options.fallbackLanguages,
    missingTranslation: options.missingTranslation,
  };
  return Object.freeze({
    localize(data, languageHeader = "") {
      return callGo({ data, languageHeader, options: configuredOptions });
    },
  });
}

class DataLocalizerError extends Error {
  constructor(message) {
    super(message);
    this.name = new.target.name;
  }
}

class ConfigurationError extends DataLocalizerError {
  constructor(message, field) {
    super(message);
    this.field = field;
  }
}

class MissingTranslationError extends DataLocalizerError {
  constructor(message, path, attemptedLanguages) {
    super(message);
    this.path = path;
    this.attemptedLanguages = attemptedLanguages;
  }
}

class CycleError extends DataLocalizerError {
  constructor(message, path, originalPath) {
    super(message);
    this.path = path;
    this.originalPath = originalPath;
  }
}

function hydrateError(details = {}) {
  switch (details.name) {
    case "ConfigurationError":
      return new ConfigurationError(details.message, details.field);
    case "MissingTranslationError":
      return new MissingTranslationError(
        details.message,
        details.path,
        details.attemptedLanguages || [],
      );
    case "CycleError":
      return new CycleError(details.message, details.path, details.originalPath);
    case "TypeError":
      return new TypeError(details.message);
    default:
      return new DataLocalizerError(details.message || "localization failed");
  }
}

exports.localize = localize;
exports.createLocalizer = createLocalizer;
exports.DataLocalizerError = DataLocalizerError;
exports.ConfigurationError = ConfigurationError;
exports.MissingTranslationError = MissingTranslationError;
exports.CycleError = CycleError;
