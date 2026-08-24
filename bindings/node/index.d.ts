export type JsonPrimitive = string | number | boolean | null;
export type JsonValue = JsonPrimitive | JsonValue[] | { [key: string]: JsonValue };
export type MissingTranslationPolicy = "preserve" | "empty" | "null" | "error";

export interface LocalizerOptions {
  supportedLanguages?: readonly string[];
  fallbackLanguage?: string;
  fallbackLanguages?: readonly string[];
  missingTranslation?: MissingTranslationPolicy;
}

export interface Localizer {
  localize<TOutput = unknown>(data: JsonValue, languageHeader?: string): TOutput;
}

export declare function localize<TOutput = unknown>(
  data: JsonValue,
  languageHeader?: string,
  fallbackLanguage?: string,
): TOutput;

export declare function createLocalizer(options?: LocalizerOptions): Localizer;

export declare class DataLocalizerError extends Error {}

export declare class ConfigurationError extends DataLocalizerError {
  readonly field?: string;
}

export declare class MissingTranslationError extends DataLocalizerError {
  readonly path?: string;
  readonly attemptedLanguages: readonly string[];
}

export declare class CycleError extends DataLocalizerError {
  readonly path?: string;
  readonly originalPath?: string;
}
