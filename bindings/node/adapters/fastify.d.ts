import type { FastifyPluginCallback, FastifyRequest } from "fastify";
import type { JsonValue } from "../index";
import type { RequestAdapterOptions } from "./framework";

export interface FastifyLocalizerOptions extends RequestAdapterOptions<FastifyRequest> {
  autoLocalizeResponse?: boolean;
  getLanguageHeader?: (request: FastifyRequest) => string | undefined;
}

declare module "fastify" {
  interface FastifyRequest {
    localizeData<TOutput = unknown>(
      data: JsonValue,
      languageHeader?: string,
    ): TOutput;
  }
}

export declare const fastifyLocalizer: FastifyPluginCallback<FastifyLocalizerOptions>;

export default fastifyLocalizer;
