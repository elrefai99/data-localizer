import { langJSON } from "../lang/lang";

export function localizeObjectFields(data: any, langHeader: string, fallbackLang: string = 'en'): any {
     const langCode = langHeader?.split(',')[0]?.toLowerCase().trim().split('-')[0];
     const preferredLang = langJSON.includes(langCode) ? langCode : fallbackLang;

     const localize = (obj: any): any => {
          if (Array.isArray(obj)) {
               return obj.map(localize);
          }

          if (obj && typeof obj === 'object') {
               if ('ar' in obj || 'en' in obj) {
                    return obj[preferredLang] || obj[fallbackLang] || '';
               }

               const result: any = {};
               for (const key in obj) {
                    result[key] = localize(obj[key]);
               }
               return result;
          }

          return obj;
     };

     return localize(data);
}
