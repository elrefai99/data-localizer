"use strict";

const { createRequestLocalizer } = require("./shared.js");

function koaLocalizer(options = {}) {
  const adapter = createRequestLocalizer(options);

  return async function dataLocalizerMiddleware(context, next) {
    context.localizeData = adapter.forRequest(context.request || context);
    await next();

    if (options.autoLocalizeResponse && context.body !== undefined) {
      context.body = context.localizeData(context.body);
    }
  };
}

exports.koaLocalizer = koaLocalizer;
exports.default = koaLocalizer;
