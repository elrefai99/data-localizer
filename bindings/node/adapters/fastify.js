"use strict";

const { createRequestLocalizer } = require("./shared.js");

function fastifyLocalizer(fastify, options = {}, done) {
  try {
    const adapter = createRequestLocalizer(options);
    fastify.decorateRequest("localizeData", function localizeData(
      data,
      languageHeader,
    ) {
      return adapter.forRequest(this)(data, languageHeader);
    });

    if (options.autoLocalizeResponse) {
      fastify.addHook(
        "preSerialization",
        function localizeResponse(request, _reply, payload, next) {
          try {
            next(null, request.localizeData(payload));
          } catch (error) {
            next(error);
          }
        },
      );
    }
    done();
  } catch (error) {
    done(error);
  }
}

module.exports = fastifyLocalizer;
module.exports.fastifyLocalizer = fastifyLocalizer;
module.exports.default = fastifyLocalizer;
