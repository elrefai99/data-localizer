import { AnyObject } from "../@types";
import { langJSON } from "../lang/lang";

function normalizeLanguage(language: string | undefined): string {
     return language?.trim().toLowerCase().split("-")[0] ?? "";
}

function localizeValue(value: unknown, preferredLanguage: string, fallbackLanguage: string,): unknown {
     if (Array.isArray(value)) {
          return value.map((item) =>
               localizeValue(item, preferredLanguage, fallbackLanguage),
          );
     }

     if (value && typeof value === "object") {
          const object = value as AnyObject;

          if ("ar" in object || "en" in object) {
               return object[preferredLanguage] ?? object[fallbackLanguage] ?? "";
          }

          const result: AnyObject = {};
          for (const key of Object.keys(object)) {
               result[key] = localizeValue(object[key], preferredLanguage, fallbackLanguage,);
          }

          return result;
     }

     return value;
}

export function localize<T extends AnyObject | AnyObject[]>(data: T, langHeader?: string, fallbackLang = "en"): T {
     const requestedLanguage = normalizeLanguage(langHeader?.split(",")[0]);
     const fallbackLanguage = normalizeLanguage(fallbackLang);
     const preferredLanguage = langJSON.includes(requestedLanguage)
          ? requestedLanguage
          : fallbackLanguage;

     return localizeValue(data, preferredLanguage, fallbackLanguage) as T;
}
