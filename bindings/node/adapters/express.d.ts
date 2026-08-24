import type { Request, RequestHandler } from "express";
import type { JsonValue } from "../index";
import type { RequestAdapterOptions } from "./framework";

export interface ExpressLocalizerOptions extends RequestAdapterOptions {
  autoLocalizeResponse?: boolean;
  getLanguageHeader?: (request: Request) => string | undefined;
}

declare global {
  namespace Express {
    interface Request {
      localizeData<TOutput = unknown>(
        data: JsonValue,
        languageHeader?: string,
      ): TOutput;
    }
  }
}

export declare function expressLocalizer(
  options?: ExpressLocalizerOptions,
): RequestHandler;

export default expressLocalizer;
