"use strict";

const { createRequestLocalizer } = require("./shared.js");

class NestDataLocalizer {
  constructor(options = {}) {
    this.adapter = createRequestLocalizer(options);
  }

  localize(data, languageHeader = "") {
    return this.adapter.localize(data, languageHeader);
  }

  forRequest(request) {
    return this.adapter.forRequest(request);
  }
}

class DataLocalizerInterceptor {
  constructor(options = {}) {
    this.adapter = createRequestLocalizer(options);
  }

  intercept(context, next) {
    const request = context.switchToHttp().getRequest();
    request.localizeData = this.adapter.forRequest(request);

    // Nest applications already depend on RxJS. Requiring it lazily keeps the
    // core and all non-Nest adapters free from this peer dependency.
    const { map } = require("rxjs/operators");
    return next
      .handle()
      .pipe(map((payload) => request.localizeData(payload)));
  }
}

function createNestInterceptor(options = {}) {
  return new DataLocalizerInterceptor(options);
}

exports.NestDataLocalizer = NestDataLocalizer;
exports.DataLocalizerInterceptor = DataLocalizerInterceptor;
exports.createNestInterceptor = createNestInterceptor;
