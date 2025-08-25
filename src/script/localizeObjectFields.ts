import { AnyObject } from "../@types";
import { langJSON } from "../lang/lang";

export function localizeObjectFields<T extends AnyObject | AnyObject[]>(data: T, langHeader: string, fallbackLang: string = 'en'): T {
     const langCode = langHeader?.split(',')[0]?.toLowerCase().trim().split('-')[0];
     const preferredLang = langJSON.includes(langCode) ? langCode : fallbackLang;

     const localize = (obj: AnyObject): AnyObject => {
          if (Array.isArray(obj)) {
               return obj.map(localize);
          }

          if (obj && typeof obj === 'object') {
               if ('ar' in obj || 'en' in obj) {
                    return obj[preferredLang] || obj[fallbackLang] || '';
               }

               const result: AnyObject = {};
               for (const key in obj) {
                    result[key] = localize(obj[key]);
               }
               return result;
          }

          return obj;
     };

     return localize(data) as T;
}
