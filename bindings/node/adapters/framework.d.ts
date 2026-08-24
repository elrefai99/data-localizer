import type { JsonValue, LocalizerOptions } from "../index";

export type RequestLike =
  | string
  | {
      headers?:
        | { get(name: string): string | null }
        | Record<string, string | readonly string[] | undefined>;
      get?(name: string): string | undefined;
    };

export interface RequestAdapterOptions extends LocalizerOptions {
  getLanguageHeader?: (request: unknown) => string | undefined;
}

export interface RequestLocalizer {
  localize<TOutput = unknown>(
    data: JsonValue,
    requestOrLanguage?: RequestLike,
  ): TOutput;
  forRequest(
    request: unknown,
  ): <TOutput = unknown>(data: JsonValue, languageHeader?: string) => TOutput;
  languageFor(request: unknown): string;
}

export declare function getAcceptLanguage(request: RequestLike): string;

export declare function createRequestLocalizer(
  options?: RequestAdapterOptions,
): RequestLocalizer;
