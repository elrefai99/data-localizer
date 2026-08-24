import assert from "node:assert/strict";

import expressDefault, {
  expressLocalizer,
} from "../../../dist/adapters/express.js";
import fastifyDefault, {
  fastifyLocalizer,
} from "../../../dist/adapters/fastify.js";
import koaDefault, { koaLocalizer } from "../../../dist/adapters/koa.js";

assert.equal(expressDefault, expressLocalizer);
assert.equal(fastifyDefault, fastifyLocalizer);
assert.equal(koaDefault, koaLocalizer);

const request = { headers: { "accept-language": "ar" } };
let localized;
expressDefault()(request, {}, () => {
  localized = request.localizeData({ title: { en: "Hello", ar: "مرحبا" } });
});
assert.deepEqual(localized, { title: "مرحبا" });

console.log("JavaScript ESM adapter smoke test passed");
