import { AnyObject } from "../@types";
import { langJSON } from "../lang/lang";
const localizeValue = (value: AnyObject, preferredLang?: any, fallbackLang?: any): AnyObject => {
     if (Array.isArray(value)) {
          return value.map(localizeValue);
     }

     if (value && typeof value === 'object') {
          if ('ar' in value || 'en' in value) {
               return value[preferredLang] ?? value[fallbackLang] ?? '';
          }

          const result: AnyObject = {};
          for (const key in value) {
               result[key] = localizeValue(value[key]);
          }
          return result;
     }

     return value;
};
export function localize<T extends AnyObject | AnyObject[]>(data: T, langHeader: string, fallbackLang: string = 'en'): T {
     const langCode = langHeader?.split(',')[0]?.toLowerCase().trim().split('-')[0] || '';
     const preferredLang = langJSON.includes(langCode) ? langCode : fallbackLang;

     return localizeValue(data, preferredLang, langCode) as T;
}
