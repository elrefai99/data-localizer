# data-localizer

A lightweight utility to localize JSON objects and arrays using language fields (e.g. `{"ar": "...", "en": "..."}`) and an HTTP language header.  
Works out of the box with **ISO 639-1** language codes (e.g., `en`, `ar`, `fr`, …) and safely falls back when a language isn’t supported.

---

## Installation

```bash
npm install data-localizer
# or
pnpm add data-localizer
```

---

## Features

- Localize **objects and arrays** recursively  
- Validates `Accept-Language` against an ISO 639-1 list (`langJSON`)  
- **Fallback** language when preferred language is missing/unsupported  
- Auto-strips region parts: `en-US` → `en`, `fr-CA` → `fr`  
- Zero dependencies, TypeScript-first  

---

## Quick Start

### Example with Arrays
```ts
import { localizeDatas } from "data-localizer";

const items = [
  { title: { ar: "مرحبا", en: "Hello" } },
  { title: { ar: "عالم",  en: "World" } }
];

localizeDatas(items, "ar");
// => [ { title: "مرحبا" }, { title: "عالم" } ]

localizeDatas(items, "en");
// => [ { title: "Hello" }, { title: "World" } ]
```

### Example with Objects
```ts
import { localizeDatas } from "data-localizer";

const user = {
  name: { ar: "محمد", en: "Mohamed" },
  age: 30
};

localizeDatas(user, "ar");
// => { name: "محمد", age: 30 }

localizeDatas(user, "en");
// => { name: "Mohamed", age: 30 }
```

### Example with Nested Structures
```ts
import { localizeObjectFields } from "data-localizer";

const payload = {
  post: {
    title: { ar: "عنوان", en: "Title" },
    meta: {
      description: { ar: "وصف", en: "Description" }
    }
  }
};

localizeObjectFields(payload, "ar");
// => { post: { title: "عنوان", meta: { description: "وصف" } } }

localizeObjectFields(payload, "en");
// => { post: { title: "Title", meta: { description: "Description" } } }
```
