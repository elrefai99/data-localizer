
export function localizeDatas<T extends any | any[]>(data: T, langHeader: string, fallbackLang: string = 'en'): T {
     const preferredLang = langHeader?.split(',')[0]?.toLowerCase().trim() === 'ar' ? 'ar' : 'en';

     const localizeObject = (obj: any): any => {
          const result: any = {};

          for (const key in obj) {
               const value = obj[key];
               if (value && typeof value === 'object' && ('ar' in value || 'en' in value)) {
                    result[key] = value[preferredLang] || value[fallbackLang] || '';
               } else {
                    result[key] = value;
               }
          }

          return result;
     };

     if (Array.isArray(data)) {
          return data.map(localizeObject) as T;
     }

     return localizeObject(data) as T;
}


export function localizeObjectFields(data: any, langHeader: string, fallbackLang: string = 'en'): any {
     const preferredLang = langHeader?.split(',')[0]?.toLowerCase().trim() === 'ar' ? 'ar' : 'en';

     const localize = (obj: any): any => {
          if (Array.isArray(obj)) {
               return obj.map(localize);
          }

          if (obj && typeof obj === 'object') {
               // Check if it's a localized value
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
