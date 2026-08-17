import { describe, expect, it } from "vitest";

import { localizeDatas, localizeObjectFields } from "../src";

const localizers = [
  ["localizeDatas", localizeDatas],
  ["localizeObjectFields", localizeObjectFields],
] as const;

describe.each(localizers)("%s", (_name, localize) => {
  it("localizes nested objects and arrays recursively", () => {
    const input = {
      title: { ar: "عنوان", en: "Title" },
      posts: [
        { name: { ar: "الأول", en: "First" }, published: true },
        { name: { ar: "الثاني", en: "Second" }, published: false },
      ],
      count: 2,
      metadata: null,
    };

    expect(localize(input, "ar")).toEqual({
      title: "عنوان",
      posts: [
        { name: "الأول", published: true },
        { name: "الثاني", published: false },
      ],
      count: 2,
      metadata: null,
    });
  });

  it("normalizes casing, whitespace, and region subtags", () => {
    const input = { greeting: { ar: "مرحبا", en: "Hello" } };

    expect(localize(input, "  AR-EG  ")).toEqual({ greeting: "مرحبا" });
  });

  it("uses the first language in an Accept-Language-style header", () => {
    const input = { greeting: { ar: "مرحبا", en: "Hello" } };

    expect(localize(input, "en-US,ar;q=0.9")).toEqual({ greeting: "Hello" });
  });

  it("uses English when the requested language is unsupported", () => {
    const input = { greeting: { ar: "مرحبا", en: "Hello" } };

    expect(localize(input, "not-a-language")).toEqual({ greeting: "Hello" });
  });

  it("supports a custom fallback language", () => {
    const input = {
      greeting: { ar: "مرحبا", en: "Hello" },
      description: { ar: "وصف", en: "Description" },
    };

    expect(localize(input, "not-a-language", "ar")).toEqual({
      greeting: "مرحبا",
      description: "وصف",
    });
  });

  it("falls back when the preferred translation is missing", () => {
    const input = { greeting: { ar: "مرحبا", en: "Hello", fr: undefined } };

    expect(localize(input, "fr", "en")).toEqual({ greeting: "Hello" });
  });

  it("does not mutate the input", () => {
    const input = {
      title: { ar: "عنوان", en: "Title" },
      nested: [{ label: { ar: "وسم", en: "Label" } }],
    };
    const snapshot = structuredClone(input);

    const result = localize(input, "en");

    expect(input).toEqual(snapshot);
    expect(result).not.toBe(input);
    expect(result.nested).not.toBe(input.nested);
  });
});

describe("nullish translation behavior", () => {
  it("localizeDatas preserves an empty preferred translation", () => {
    expect(localizeDatas({ title: { ar: "", en: "Title" } }, "ar")).toEqual({
      title: "",
    });
  });

  it("localizeObjectFields falls back from an empty preferred translation", () => {
    expect(
      localizeObjectFields({ title: { ar: "", en: "Title" } }, "ar"),
    ).toEqual({ title: "Title" });
  });
});
