"use strict";

const { createRequestLocalizer } = require("./shared.js");

function expressLocalizer(options = {}) {
  const adapter = createRequestLocalizer(options);

  return function dataLocalizerMiddleware(request, response, next) {
    request.localizeData = adapter.forRequest(request);

    if (options.autoLocalizeResponse && typeof response.json === "function") {
      const sendJSON = response.json.bind(response);
      response.json = (body) => sendJSON(request.localizeData(body));
    }

    next();
  };
}

exports.expressLocalizer = expressLocalizer;
exports.default = expressLocalizer;
