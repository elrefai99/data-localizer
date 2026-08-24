"use strict";

const { createLocalizer } = require("../index.js");

function getAcceptLanguage(request) {
  if (typeof request === "string") {
    return request;
  }
  if (!request || typeof request !== "object") {
    return "";
  }

  if (typeof request.get === "function") {
    const value = request.get("Accept-Language");
    if (typeof value === "string") {
      return value;
    }
  }

  const headers = request.headers;
  if (!headers) {
    return "";
  }
  if (typeof headers.get === "function") {
    return headers.get("accept-language") || "";
  }

  const key = Object.keys(headers).find(
    (name) => name.toLowerCase() === "accept-language",
  );
  const value = key ? headers[key] : undefined;
  return Array.isArray(value) ? value.join(",") : value || "";
}

function toLocalizerOptions(options) {
  return {
    supportedLanguages: options.supportedLanguages,
    fallbackLanguage: options.fallbackLanguage,
    fallbackLanguages: options.fallbackLanguages,
    missingTranslation: options.missingTranslation,
  };
}

function createRequestLocalizer(options = {}) {
  const engine = createLocalizer(toLocalizerOptions(options));
  const readLanguage = options.getLanguageHeader || getAcceptLanguage;

  function languageFor(request) {
    const value = readLanguage(request);
    return typeof value === "string" ? value : "";
  }

  return Object.freeze({
    localize(data, requestOrLanguage = "") {
      const language =
        typeof requestOrLanguage === "string"
          ? requestOrLanguage
          : languageFor(requestOrLanguage);
      return engine.localize(data, language);
    },
    forRequest(request) {
      return (data, languageHeader) =>
        engine.localize(
          data,
          languageHeader === undefined ? languageFor(request) : languageHeader,
        );
    },
    languageFor,
  });
}

exports.getAcceptLanguage = getAcceptLanguage;
exports.createRequestLocalizer = createRequestLocalizer;
