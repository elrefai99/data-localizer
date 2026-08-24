# data-localizer

Localize language-keyed JSON data in Node.js applications using an
`Accept-Language` value. The package supports TypeScript, JavaScript ESM, and
CommonJS on Node.js 18 or newer.

It performs no network requests and has no required npm dependencies.

## Installation

```bash
pnpm install data-localizer
```

## TypeScript and JavaScript ESM

```ts
import { localize } from "data-localizer";

const product = {
  id: 42,
  title: {
    en: "Coffee",
    ar: "قهوة",
  },
};

const result = localize(product, "ar-EG,en;q=0.8");
// { id: 42, title: "قهوة" }
```

## CommonJS JavaScript

```js
const { localize } = require("data-localizer");

const result = localize(
  { title: { en: "Coffee", ar: "قهوة" } },
  "ar",
);
```

## Configuration

Create a reusable localizer for an application with a known language set:

```ts
import { createLocalizer } from "data-localizer";

const localizer = createLocalizer({
  supportedLanguages: ["en", "ar", "fr"],
  fallbackLanguage: "en",
  fallbackLanguages: ["fr"],
  missingTranslation: "preserve",
});

const result = localizer.localize(data, "fr-CA,fr;q=0.9,en;q=0.8");
```

Missing-translation policies:

- `preserve`: keep the original translation object.
- `empty`: return an empty string.
- `null`: return `null`.
- `error`: throw `MissingTranslationError`.

Empty strings, `false`, and zero are treated as valid translations. Input data
is not mutated.

## Express

Install Express normally, then register the adapter:

```bash
pnpm install express data-localizer
```

TypeScript or JavaScript ESM:

```ts
import express from "express";
import expressLocalizer from "data-localizer/express";

const app = express();

app.use(expressLocalizer({
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
}));

app.get("/product", (request, response) => {
  response.json(request.localizeData({
    title: { en: "Coffee", ar: "قهوة" },
  }));
});
```

CommonJS:

```js
const express = require("express");
const expressLocalizer = require("data-localizer/express");

const app = express();
app.use(expressLocalizer({ fallbackLanguage: "en" }));
```

Automatic response localization is also available:

```ts
app.use(expressLocalizer({
  fallbackLanguage: "en",
  autoLocalizeResponse: true,
}));

app.get("/product", (_request, response) => {
  response.json({ title: { en: "Coffee", ar: "قهوة" } });
});
```

## NestJS

```ts
import { createNestInterceptor } from "data-localizer/nest";

app.useGlobalInterceptors(createNestInterceptor({
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
}));
```

`NestDataLocalizer` is also available for registration as a custom provider or
for direct use in a service.

## Fastify

```ts
import fastifyLocalizer from "data-localizer/fastify";

await app.register(fastifyLocalizer, {
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
  autoLocalizeResponse: true,
});
```

Every request receives `request.localizeData(data)`.

## Koa

```ts
import koaLocalizer from "data-localizer/koa";

app.use(koaLocalizer({
  fallbackLanguage: "en",
  autoLocalizeResponse: true,
}));
```

Every context receives `context.localizeData(data)`.

## Other frameworks

Use the framework-neutral request adapter:

```ts
import { createRequestLocalizer } from "data-localizer/framework";

const adapter = createRequestLocalizer({ fallbackLanguage: "en" });
const result = adapter.localize(data, request);
```

By default, the adapter reads `Accept-Language` from `request.headers`. A custom
reader can be supplied when a framework stores it elsewhere:

```ts
const adapter = createRequestLocalizer<{ locale?: string }>({
  getLanguageHeader: (request) => request.locale,
});
```

## License

MIT
