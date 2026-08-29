"use strict";

const assert = require("node:assert/strict");
const {
  MissingTranslationError,
  createLocalizer,
  localize,
} = require("../../dist/index.js");

const input = {
  title: { en: "Hello", ar: "\u0645\u0631\u062d\u0628\u0627" },
  items: [{ name: { en: "First", ar: "\u0627\u0644\u0623\u0648\u0644" } }],
};

assert.deepEqual(localize(input, "ar-EG,en;q=0.8"), {
  title: "\u0645\u0631\u062d\u0628\u0627",
  items: [{ name: "\u0627\u0644\u0623\u0648\u0644" }],
});
assert.deepEqual(input.title, { en: "Hello", ar: "\u0645\u0631\u062d\u0628\u0627" });

const strict = createLocalizer({
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
  missingTranslation: "error",
});
assert.throws(
  () => strict.localize({ title: { ar: null } }, "ar"),
  (error) =>
    error instanceof MissingTranslationError &&
    error.path === "$.title" &&
    error.attemptedLanguages.join(",") === "ar,en",
);

console.log("npm WebAssembly smoke test passed");
