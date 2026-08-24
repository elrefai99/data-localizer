import type {
  CallHandler,
  ExecutionContext,
  NestInterceptor,
} from "@nestjs/common";
import type { Observable } from "rxjs";
import type { JsonValue } from "../index";
import type { RequestAdapterOptions } from "./framework";

export type RequestLocalizeFunction = <TOutput = unknown>(
  data: JsonValue,
  languageHeader?: string,
) => TOutput;

export declare class NestDataLocalizer {
  constructor(options?: RequestAdapterOptions<unknown>);
  localize<TOutput = unknown>(
    data: JsonValue,
    languageHeader?: string,
  ): TOutput;
  forRequest(request: unknown): RequestLocalizeFunction;
}

export declare class DataLocalizerInterceptor implements NestInterceptor {
  constructor(options?: RequestAdapterOptions<unknown>);
  intercept(context: ExecutionContext, next: CallHandler): Observable<unknown>;
}

export declare function createNestInterceptor(
  options?: RequestAdapterOptions<unknown>,
): DataLocalizerInterceptor;
