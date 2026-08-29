import assert from "node:assert/strict";

import expressDefault, {
  expressLocalizer,
} from "../../dist/adapters/express.js";
import fastifyDefault, {
  fastifyLocalizer,
} from "../../dist/adapters/fastify.js";
import koaDefault, { koaLocalizer } from "../../dist/adapters/koa.js";

assert.equal(expressDefault, expressLocalizer);
assert.equal(fastifyDefault, fastifyLocalizer);
assert.equal(koaDefault, koaLocalizer);

const request = { headers: { "accept-language": "ar" } };
let localized;
expressDefault()(request, {}, () => {
  localized = request.localizeData({
    title: { en: "Hello", ar: "\u0645\u0631\u062d\u0628\u0627" },
  });
});
assert.deepEqual(localized, { title: "\u0645\u0631\u062d\u0628\u0627" });

console.log("JavaScript ESM adapter smoke test passed");
