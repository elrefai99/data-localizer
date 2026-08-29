package main

import "strings"

type tsParameter struct {
	name     string
	typeName string
	optional bool
}

type tsProperty struct {
	name     string
	typeName string
	optional bool
	readonly bool
}

type tsMethod struct {
	name           string
	typeParameters string
	parameters     []tsParameter
	returnType     string
}

type tsInterface struct {
	name       string
	extends    string
	properties []tsProperty
	methods    []tsMethod
}

type tsFunction struct {
	name           string
	typeParameters string
	parameters     []tsParameter
	returnType     string
}

type tsClass struct {
	name        string
	extends     string
	implements  string
	constructor []tsParameter
	properties  []tsProperty
	methods     []tsMethod
}

type declarationWriter struct {
	sourceWriter
}

func declarationFiles() map[string]string {
	return map[string]string{
		"index.d.ts":              indexDeclarations(),
		"adapters/framework.d.ts": frameworkDeclarations(),
		"adapters/express.d.ts":   expressDeclarations(),
		"adapters/nest.d.ts":      nestDeclarations(),
		"adapters/fastify.d.ts":   fastifyDeclarations(),
		"adapters/koa.d.ts":       koaDeclarations(),
	}
}

func newDeclarationWriter() *declarationWriter {
	w := &declarationWriter{}
	w.WriteString(generatedNotice)
	return w
}

func (w *declarationWriter) typeImport(module string, names ...string) {
	w.line(`import type { %s } from %q;`, strings.Join(names, ", "), module)
}

func (w *declarationWriter) defaultImport(name, module string) {
	w.line(`import %s = require(%q);`, name, module)
}

func (w *declarationWriter) alias(name, value string) {
	w.line(`export type %s = %s;`, name, value)
}

func (w *declarationWriter) interfaceType(value tsInterface, indent string, exported bool) {
	prefix := ""
	if exported {
		prefix = "export "
	}
	extends := ""
	if value.extends != "" {
		extends = " extends " + value.extends
	}
	w.line(`%s%sinterface %s%s {`, indent, prefix, value.name, extends)
	for _, property := range value.properties {
		w.property(property, indent+"  ")
	}
	for _, method := range value.methods {
		w.method(method, indent+"  ")
	}
	w.line(`%s}`, indent)
}

func (w *declarationWriter) function(value tsFunction) {
	w.line(`export declare function %s%s(%s): %s;`, value.name, value.typeParameters, parameters(value.parameters), value.returnType)
}

func (w *declarationWriter) class(value tsClass) {
	extends := ""
	if value.extends != "" {
		extends = " extends " + value.extends
	}
	implements := ""
	if value.implements != "" {
		implements = " implements " + value.implements
	}
	w.line(`export declare class %s%s%s {`, value.name, extends, implements)
	if value.constructor != nil {
		w.line(`  constructor(%s);`, parameters(value.constructor))
	}
	for _, property := range value.properties {
		w.property(property, "  ")
	}
	for _, method := range value.methods {
		w.method(method, "  ")
	}
	w.line(`}`)
}

func (w *declarationWriter) property(value tsProperty, indent string) {
	readonly := ""
	if value.readonly {
		readonly = "readonly "
	}
	optional := ""
	if value.optional {
		optional = "?"
	}
	w.line(`%s%s%s%s: %s;`, indent, readonly, value.name, optional, value.typeName)
}

func (w *declarationWriter) method(value tsMethod, indent string) {
	w.line(`%s%s%s(%s): %s;`, indent, value.name, value.typeParameters, parameters(value.parameters), value.returnType)
}

func parameters(values []tsParameter) string {
	result := make([]string, len(values))
	for index, value := range values {
		optional := ""
		if value.optional {
			optional = "?"
		}
		result[index] = value.name + optional + ": " + value.typeName
	}
	return strings.Join(result, ", ")
}

func indexDeclarations() string {
	w := newDeclarationWriter()
	w.alias("JsonPrimitive", "string | number | boolean | null")
	w.alias("JsonValue", "JsonPrimitive | JsonValue[] | { [key: string]: JsonValue }")
	w.alias("MissingTranslationPolicy", `"preserve" | "empty" | "null" | "error"`)
	w.line("")
	w.interfaceType(tsInterface{name: "LocalizerOptions", properties: []tsProperty{
		{name: "supportedLanguages", typeName: "readonly string[]", optional: true},
		{name: "fallbackLanguage", typeName: "string", optional: true},
		{name: "fallbackLanguages", typeName: "readonly string[]", optional: true},
		{name: "missingTranslation", typeName: "MissingTranslationPolicy", optional: true},
	}}, "", true)
	w.line("")
	w.interfaceType(tsInterface{name: "Localizer", methods: []tsMethod{localizeMethod()}}, "", true)
	w.line("")
	w.function(tsFunction{name: "localize", typeParameters: "<TOutput = unknown>", parameters: []tsParameter{
		{name: "data", typeName: "JsonValue"},
		{name: "languageHeader", typeName: "string", optional: true},
		{name: "fallbackLanguage", typeName: "string", optional: true},
	}, returnType: "TOutput"})
	w.function(tsFunction{name: "createLocalizer", parameters: []tsParameter{{name: "options", typeName: "LocalizerOptions", optional: true}}, returnType: "Localizer"})
	w.line("")
	w.class(tsClass{name: "DataLocalizerError", extends: "Error"})
	w.class(tsClass{name: "ConfigurationError", extends: "DataLocalizerError", properties: []tsProperty{{name: "field", typeName: "string", optional: true, readonly: true}}})
	w.class(tsClass{name: "MissingTranslationError", extends: "DataLocalizerError", properties: []tsProperty{
		{name: "path", typeName: "string", optional: true, readonly: true},
		{name: "attemptedLanguages", typeName: "readonly string[]", readonly: true},
	}})
	w.class(tsClass{name: "CycleError", extends: "DataLocalizerError", properties: []tsProperty{
		{name: "path", typeName: "string", optional: true, readonly: true},
		{name: "originalPath", typeName: "string", optional: true, readonly: true},
	}})
	return w.String()
}

func frameworkDeclarations() string {
	w := newDeclarationWriter()
	w.typeImport("../index", "JsonValue", "LocalizerOptions")
	w.line("")
	w.alias("RequestLike", `string | { headers?: { get(name: string): string | null } | Record<string, string | readonly string[] | undefined>; get?(name: string): string | undefined }`)
	w.line("")
	w.interfaceType(tsInterface{name: "RequestAdapterOptions<TRequest = unknown>", extends: "LocalizerOptions", properties: []tsProperty{
		{name: "getLanguageHeader", typeName: "(request: TRequest) => string | undefined", optional: true},
	}}, "", true)
	w.line("")
	w.interfaceType(tsInterface{name: "RequestLocalizer<TRequest = RequestLike>", methods: []tsMethod{
		{name: "localize", typeParameters: "<TOutput = unknown>", parameters: []tsParameter{{name: "data", typeName: "JsonValue"}, {name: "requestOrLanguage", typeName: "TRequest | string", optional: true}}, returnType: "TOutput"},
		{name: "forRequest", parameters: []tsParameter{{name: "request", typeName: "TRequest"}}, returnType: "<TOutput = unknown>(data: JsonValue, languageHeader?: string) => TOutput"},
		{name: "languageFor", parameters: []tsParameter{{name: "request", typeName: "TRequest"}}, returnType: "string"},
	}}, "", true)
	w.line("")
	w.function(tsFunction{name: "getAcceptLanguage", parameters: []tsParameter{{name: "request", typeName: "RequestLike"}}, returnType: "string"})
	w.function(tsFunction{name: "createRequestLocalizer", typeParameters: "<TRequest = RequestLike>", parameters: []tsParameter{{name: "options", typeName: "RequestAdapterOptions<TRequest>", optional: true}}, returnType: "RequestLocalizer<TRequest>"})
	return w.String()
}

func expressDeclarations() string {
	w := newDeclarationWriter()
	w.typeImport("express", "Request", "RequestHandler")
	w.typeImport("../index", "JsonValue")
	w.typeImport("./framework", "RequestAdapterOptions")
	w.line("")
	w.interfaceType(tsInterface{name: "ExpressLocalizerOptions", extends: "RequestAdapterOptions<Request>", properties: []tsProperty{
		{name: "autoLocalizeResponse", typeName: "boolean", optional: true},
		{name: "getLanguageHeader", typeName: "(request: Request) => string | undefined", optional: true},
	}}, "", true)
	w.line("")
	w.line(`declare global {`)
	w.line(`  namespace Express {`)
	w.interfaceType(tsInterface{name: "Request", methods: []tsMethod{localizeMethod()}}, "    ", false)
	w.line(`  }`)
	w.line(`}`)
	w.line("")
	w.function(tsFunction{name: "expressLocalizer", parameters: []tsParameter{{name: "options", typeName: "ExpressLocalizerOptions", optional: true}}, returnType: "RequestHandler"})
	w.line(`export default expressLocalizer;`)
	return w.String()
}

func nestDeclarations() string {
	w := newDeclarationWriter()
	w.typeImport("@nestjs/common", "CallHandler", "ExecutionContext", "NestInterceptor")
	w.typeImport("rxjs", "Observable")
	w.typeImport("../index", "JsonValue")
	w.typeImport("./framework", "RequestAdapterOptions")
	w.line("")
	w.alias("RequestLocalizeFunction", "<TOutput = unknown>(data: JsonValue, languageHeader?: string) => TOutput")
	w.line("")
	w.class(tsClass{name: "NestDataLocalizer", constructor: []tsParameter{{name: "options", typeName: "RequestAdapterOptions<unknown>", optional: true}}, methods: []tsMethod{
		localizeMethod(),
		{name: "forRequest", parameters: []tsParameter{{name: "request", typeName: "unknown"}}, returnType: "RequestLocalizeFunction"},
	}})
	w.class(tsClass{name: "DataLocalizerInterceptor", implements: "NestInterceptor", constructor: []tsParameter{{name: "options", typeName: "RequestAdapterOptions<unknown>", optional: true}}, methods: []tsMethod{
		{name: "intercept", parameters: []tsParameter{{name: "context", typeName: "ExecutionContext"}, {name: "next", typeName: "CallHandler"}}, returnType: "Observable<unknown>"},
	}})
	w.function(tsFunction{name: "createNestInterceptor", parameters: []tsParameter{{name: "options", typeName: "RequestAdapterOptions<unknown>", optional: true}}, returnType: "DataLocalizerInterceptor"})
	return w.String()
}

func fastifyDeclarations() string {
	w := newDeclarationWriter()
	w.typeImport("fastify", "FastifyPluginCallback", "FastifyRequest")
	w.typeImport("../index", "JsonValue")
	w.typeImport("./framework", "RequestAdapterOptions")
	w.line("")
	w.interfaceType(tsInterface{name: "FastifyLocalizerOptions", extends: "RequestAdapterOptions<FastifyRequest>", properties: []tsProperty{
		{name: "autoLocalizeResponse", typeName: "boolean", optional: true},
		{name: "getLanguageHeader", typeName: "(request: FastifyRequest) => string | undefined", optional: true},
	}}, "", true)
	w.line("")
	w.line(`declare module "fastify" {`)
	w.interfaceType(tsInterface{name: "FastifyRequest", methods: []tsMethod{localizeMethod()}}, "  ", false)
	w.line(`}`)
	w.line("")
	w.line(`export declare const fastifyLocalizer: FastifyPluginCallback<FastifyLocalizerOptions>;`)
	w.line(`export default fastifyLocalizer;`)
	return w.String()
}

func koaDeclarations() string {
	w := newDeclarationWriter()
	w.defaultImport("Koa", "koa")
	w.typeImport("../index", "JsonValue")
	w.typeImport("./framework", "RequestAdapterOptions")
	w.line("")
	w.interfaceType(tsInterface{name: "KoaLocalizerOptions", extends: "RequestAdapterOptions<Koa.Request>", properties: []tsProperty{
		{name: "autoLocalizeResponse", typeName: "boolean", optional: true},
		{name: "getLanguageHeader", typeName: "(request: Koa.Request) => string | undefined", optional: true},
	}}, "", true)
	w.line("")
	w.line(`declare module "koa" {`)
	w.interfaceType(tsInterface{name: "BaseContext", methods: []tsMethod{localizeMethod()}}, "  ", false)
	w.line(`}`)
	w.line("")
	w.function(tsFunction{name: "koaLocalizer", parameters: []tsParameter{{name: "options", typeName: "KoaLocalizerOptions", optional: true}}, returnType: "Koa.Middleware"})
	w.line(`export default koaLocalizer;`)
	return w.String()
}

func localizeMethod() tsMethod {
	return tsMethod{name: "localize", typeParameters: "<TOutput = unknown>", parameters: []tsParameter{
		{name: "data", typeName: "JsonValue"},
		{name: "languageHeader", typeName: "string", optional: true},
	}, returnType: "TOutput"}
}
