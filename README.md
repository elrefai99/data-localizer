# 📦 data-localizer

A lightweight utility to localize JSON objects and arrays between Arabic (`ar`) and English (`en`) fields based on a language header.  
Perfect for APIs and backends that store multilingual data inside objects like:

```json
{
  "title": { "ar": "مرحبا", "en": "Hello" }
}
```

---

## 🚀 Installation

Install using **npm**:

```bash
npm install data-localizer
```

Or using **pnpm**:

```bash
pnpm add data-localizer
```

---

## 🛠 Usage

### Example with Array

```ts
import { localizeDatas } from "data-localizer";

const data = [
  { title: { ar: "مرحبا", en: "Hello" } },
  { title: { ar: "عالم", en: "World" } }
];

console.log(localizeDatas(data, "ar"));
// => [ { title: "مرحبا" }, { title: "عالم" } ]

console.log(localizeDatas(data, "en"));
// => [ { title: "Hello" }, { title: "World" } ]
```

---

### Example with Single Object

```ts
import { localizeDatas } from "data-localizer";

const user = {
  name: { ar: "محمد", en: "Mohamed" },
  age: 28
};

console.log(localizeDatas(user, "ar"));
// => { name: "محمد", age: 28 }

console.log(localizeDatas(user, "en"));
// => { name: "Mohamed", age: 28 }
```

---

## ⚙️ API

### `localizeDatas<T>(data: T, langHeader: string, fallbackLang?: string): T`

#### Parameters:
- **`data`**: `object | object[]`  
  The data you want to localize (can be a single object or an array of objects).  

- **`langHeader`**: `string`  
  The language header (usually from HTTP request headers like `Accept-Language`).  
  - If `langHeader` starts with `"ar"`, Arabic will be used.  
  - Otherwise, English will be used.  

- **`fallbackLang`** *(optional)*: `string` (default: `"en"`)  
  The fallback language in case the preferred language value is missing.  

#### Returns:
- A new object or array of objects with localized values.

---

## 🌍 Why use data-localizer?

- 🔑 **Simple & Lightweight** → No heavy dependencies.  
- 🌐 **Multi-language ready** → Works out of the box for `ar/en`, and extendable for other languages.  
- ⚡ **API-friendly** → Great for backends that serve localized content.  
- 💡 **TypeScript support**
