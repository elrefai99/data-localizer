# data-localizer

`data-localizer` is a dependency-free Go library and CLI for recursively
localizing JSON-compatible data. It resolves objects such as
`{"en":"Hello","ar":"مرحبا"}` using an `Accept-Language` value.

It performs no network requests and does not call any external API.

The repository also builds an npm package. Its localization logic is the same
Go code compiled to WebAssembly; JavaScript only provides the Node.js wrapper.

## Install

```bash
go get github.com/elrefai99/data-localizer
```

The module requires Go 1.22 or newer and uses only the standard library.

## Library usage

```go
package main

import (
	"fmt"
	"log"

	localizer "github.com/elrefai99/data-localizer"
)

func main() {
	data := map[string]any{
		"title": map[string]any{
			"en": "Coffee",
			"ar": "قهوة",
		},
	}

	result, err := localizer.Localize(data, "ar-EG, en;q=0.8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", result)
	// map[title:قهوة]
}
```

`Localize` uses English as its fallback. A third argument changes it:

```go
result, err := localizer.Localize(data, "unknown", "ar")
```

## Configured localizer

Create one reusable localizer when the application has a known language set or
needs an explicit missing-value policy:

```go
engine, err := localizer.New(localizer.Options{
	SupportedLanguages: []string{"en", "ar", "fr", "fr-CA"},
	FallbackLanguage:   "en",
	FallbackLanguages:  []string{"fr"},
	MissingTranslation: localizer.MissingPreserve,
})
if err != nil {
	log.Fatal(err)
}

result, err := engine.Localize(data, "fr-CA,fr;q=0.9,en;q=0.8")
```

Available policies are `MissingPreserve`, `MissingEmpty`, `MissingNull`, and
`MissingError`. Custom detection and missing handlers are available through
`Options.IsTranslationMap` and `Options.OnMissing`.

The engine:

- parses quality values and regional language tags;
- keeps valid empty strings, `false`, and zero translations;
- recursively copies maps, slices, and arrays without mutating the input;
- treats only supported-language scalar maps as translations by default;
- preserves structs and other non-JSON values unchanged;
- returns `CycleError` for cyclic maps or slices.

## CLI

Build the command:

```bash
go build -o bin/data-localizer ./cmd/data-localizer
```

Read JSON from stdin:

```bash
echo '{"title":{"en":"Hello","ar":"مرحبا"}}' |
  ./bin/data-localizer -lang 'ar-EG, en;q=0.8' -pretty
```

Or pass a file:

```bash
./bin/data-localizer -lang fr -fallback en input.json
```

CLI flags:

- `-lang`: an `Accept-Language` value;
- `-fallback`: the primary fallback language (default `en`);
- `-supported`: optional comma-separated language list;
- `-missing`: `preserve`, `empty`, `null`, or `error`;
- `-pretty`: pretty-print output JSON.

## npm package

Build the JavaScript wrapper, TypeScript declarations, and Go WebAssembly
binary:

```bash
npm run build
```

Use it from JavaScript or TypeScript with the familiar synchronous API:

```ts
import { localize, createLocalizer } from "data-localizer";

const result = localize(
  { title: { en: "Coffee", ar: "قهوة" } },
  "ar-EG,en;q=0.8",
);
// { title: "قهوة" }

const engine = createLocalizer({
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
  missingTranslation: "preserve",
});
```

Create and inspect the archive before publishing:

```bash
npm pack --dry-run
npm publish
```

The published package has no npm dependencies. It targets Node.js 18 or newer
and contains `index.js`, `index.d.ts`, the Go runtime, and the compiled `.wasm`
engine.

## Framework adapters

Frameworks are optional peer dependencies. Install only the framework used by
the application.

### Express

```ts
import express from "express";
import { expressLocalizer } from "data-localizer/express";

const app = express();
app.use(expressLocalizer({
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
  autoLocalizeResponse: true,
}));

app.get("/product", (_request, response) => {
  response.json({ title: { en: "Coffee", ar: "قهوة" } });
});
```

Without `autoLocalizeResponse`, call `request.localizeData(data)` explicitly.

CommonJS JavaScript is also supported:

```js
const express = require("express");
const expressLocalizer = require("data-localizer/express");

const app = express();
app.use(expressLocalizer({ fallbackLanguage: "en" }));

app.get("/product", (request, response) => {
  response.json(request.localizeData({
    title: { en: "Coffee", ar: "قهوة" },
  }));
});
```

For JavaScript ESM, use the same default import shown in the TypeScript
example. Named imports also work:

```js
import { expressLocalizer } from "data-localizer/express";
```

### NestJS

```ts
import { createNestInterceptor } from "data-localizer/nest";

app.useGlobalInterceptors(createNestInterceptor({
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
}));
```

`NestDataLocalizer` is also exported for use as a service or custom provider.

### Fastify

```ts
import fastifyLocalizer from "data-localizer/fastify";

await app.register(fastifyLocalizer, {
  supportedLanguages: ["en", "ar"],
  fallbackLanguage: "en",
  autoLocalizeResponse: true,
});
```

Every Fastify request receives `request.localizeData(data)`.

### Koa

```ts
import { koaLocalizer } from "data-localizer/koa";

app.use(koaLocalizer({
  fallbackLanguage: "en",
  autoLocalizeResponse: true,
}));
```

For other Node.js or TypeScript frameworks, use the neutral adapter:

```ts
import { createRequestLocalizer } from "data-localizer/framework";

const adapter = createRequestLocalizer({ fallbackLanguage: "en" });
const result = adapter.localize(data, request);
```

## Development

```bash
go test ./...
go vet ./...
go build ./cmd/data-localizer
```

## License

MIT
