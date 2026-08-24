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

module.exports = koaLocalizer;
module.exports.koaLocalizer = koaLocalizer;
module.exports.default = koaLocalizer;
