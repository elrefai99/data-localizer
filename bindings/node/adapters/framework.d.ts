import type { JsonValue, LocalizerOptions } from "../index";

export type RequestLike =
  | string
  | {
      headers?:
        | { get(name: string): string | null }
        | Record<string, string | readonly string[] | undefined>;
      get?(name: string): string | undefined;
    };

export interface RequestAdapterOptions<TRequest = unknown> extends LocalizerOptions {
  getLanguageHeader?: (request: TRequest) => string | undefined;
}

export interface RequestLocalizer<TRequest = RequestLike> {
  localize<TOutput = unknown>(
    data: JsonValue,
    requestOrLanguage?: TRequest | string,
  ): TOutput;
  forRequest(
    request: TRequest,
  ): <TOutput = unknown>(data: JsonValue, languageHeader?: string) => TOutput;
  languageFor(request: TRequest): string;
}

export declare function getAcceptLanguage(request: RequestLike): string;

export declare function createRequestLocalizer<TRequest = RequestLike>(
  options?: RequestAdapterOptions<TRequest>,
): RequestLocalizer<TRequest>;
