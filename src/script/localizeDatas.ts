import { langJSON } from "../lang/lang";

export function localizeDatas<T extends any | any[]>(data: T, langHeader: string, fallbackLang: string = 'en'): T {
     const langCode = langHeader?.split(',')[0]?.toLowerCase().trim().split('-')[0] || '';
     const preferredLang = langJSON.includes(langCode) ? langCode : fallbackLang;

     const localizeValue = (value: any): any => {
          if (Array.isArray(value)) {
               return value.map(localizeValue);
          }

          if (value && typeof value === 'object') {
               if ('ar' in value || 'en' in value) {
                    return value[preferredLang] ?? value[fallbackLang] ?? '';
               }

               const result: any = {};
               for (const key in value) {
                    result[key] = localizeValue(value[key]);
               }
               return result;
          }

          return value;
     };

     return localizeValue(data) as T;
}
