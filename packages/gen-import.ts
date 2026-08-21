// @ts-nocheck — auto-generated barrel with lazy CJS re-exports
/**
 * gen-import.ts — AUTO-GENERATED, do not edit manually.
 * Regenerate: npx gen-import
 *
 * Value exports use lazy getters to prevent circular-dependency
 * errors when source files import from this barrel (CJS).
 *
 * Getters are installed on module.exports, not exports: esbuild-based
 * loaders (tsx, bun) reassign module.exports for any file containing
 * export syntax, which would strand getters bound to exports.
 */

export declare const langJSON: typeof import('./lang/lang').langJSON;
export declare const localize: typeof import('./localizeDatas/localizeDatas').localize;

Object.defineProperty(module.exports, 'langJSON', { get() { return require('./lang/lang').langJSON }, enumerable: true, configurable: true });
Object.defineProperty(module.exports, 'localize', { get() { return require('./localizeDatas/localizeDatas').localize }, enumerable: true, configurable: true });
