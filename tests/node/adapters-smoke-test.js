"use strict";

const assert = require("node:assert/strict");
const {
  createRequestLocalizer,
  getAcceptLanguage,
} = require("../../dist/adapters/framework.js");
const expressModule = require("../../dist/adapters/express.js");
const { expressLocalizer } = expressModule;
const { NestDataLocalizer } = require("../../dist/adapters/nest.js");
const { fastifyLocalizer } = require("../../dist/adapters/fastify.js");
const { koaLocalizer } = require("../../dist/adapters/koa.js");

const data = { title: { en: "Hello", ar: "\u0645\u0631\u062d\u0628\u0627" } };

async function main() {
  assert.equal(typeof expressModule, "function");
  assert.equal(expressModule, expressLocalizer);
  assert.equal(
    getAcceptLanguage({ headers: { "Accept-Language": ["ar", "en;q=0.8"] } }),
    "ar,en;q=0.8",
  );

  const generic = createRequestLocalizer();
  assert.deepEqual(
    generic.localize(data, { headers: { "accept-language": "ar-EG" } }),
    { title: "\u0645\u0631\u062d\u0628\u0627" },
  );

  let expressOutput;
  const expressRequest = { headers: { "accept-language": "ar" } };
  const expressResponse = {
    json(body) {
      expressOutput = body;
      return this;
    },
  };
  expressLocalizer({ autoLocalizeResponse: true })(
    expressRequest,
    expressResponse,
    () => expressResponse.json(data),
  );
  assert.deepEqual(expressOutput, { title: "\u0645\u0631\u062d\u0628\u0627" });
  assert.deepEqual(expressRequest.localizeData(data, "en"), { title: "Hello" });

  const nest = new NestDataLocalizer();
  assert.deepEqual(
    nest.forRequest({ headers: { "accept-language": "ar" } })(data),
    { title: "\u0645\u0631\u062d\u0628\u0627" },
  );

  let fastifyDecorator;
  let fastifyHook;
  fastifyLocalizer(
    {
      decorateRequest(_name, value) {
        fastifyDecorator = value;
      },
      addHook(_name, hook) {
        fastifyHook = hook;
      },
    },
    { autoLocalizeResponse: true },
    (error) => {
      if (error) throw error;
    },
  );
  const fastifyRequest = { headers: { "accept-language": "ar" } };
  fastifyRequest.localizeData = fastifyDecorator;
  await new Promise((resolve, reject) => {
    fastifyHook(fastifyRequest, {}, data, (error, result) => {
      if (error) return reject(error);
      assert.deepEqual(result, { title: "\u0645\u0631\u062d\u0628\u0627" });
      resolve();
    });
  });

  const koaContext = {
    request: { headers: { "accept-language": "ar" } },
    body: data,
  };
  await koaLocalizer({ autoLocalizeResponse: true })(koaContext, async () => {});
  assert.deepEqual(koaContext.body, { title: "\u0645\u0631\u062d\u0628\u0627" });

  console.log("framework adapter smoke tests passed");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
