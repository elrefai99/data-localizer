"use strict";

const assert = require("node:assert/strict");
const {
  MissingTranslationError,
  createLocalizer,
  localize,
} = require("../../dist/index.js");

const input = {
  title: { en: "Hello", ar: "مرحبا" },
  items: [{ name: { en: "First", ar: "الأول" } }],
};

assert.deepEqual(localize(input, "ar-EG,en;q=0.8"), {
  title: "مرحبا",
  items: [{ name: "الأول" }],
});
assert.deepEqual(input.title, { en: "Hello", ar: "مرحبا" });

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
