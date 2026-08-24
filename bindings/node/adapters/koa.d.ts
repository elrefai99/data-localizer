import Koa = require("koa");
import type { JsonValue } from "../index";
import type { RequestAdapterOptions } from "./framework";

export interface KoaLocalizerOptions extends RequestAdapterOptions<Koa.Request> {
  autoLocalizeResponse?: boolean;
  getLanguageHeader?: (request: Koa.Request) => string | undefined;
}

declare module "koa" {
  interface BaseContext {
    localizeData<TOutput = unknown>(
      data: JsonValue,
      languageHeader?: string,
    ): TOutput;
  }
}

export declare function koaLocalizer(options?: KoaLocalizerOptions): Koa.Middleware;

export default koaLocalizer;
