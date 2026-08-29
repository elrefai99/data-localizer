//go:build js && wasm

package main

import "syscall/js"

func (b *bridge) frameworkModule() js.Value {
	module := object()
	module.Set("getAcceptLanguage", b.function(func(_ js.Value, arguments []js.Value) any {
		return acceptLanguage(valueArgument(arguments, 0))
	}))
	module.Set("createRequestLocalizer", b.function(b.createRequestLocalizer))
	return module
}

func (b *bridge) createRequestLocalizer(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		adapter, err := b.newRequestLocalizer(valueArgument(arguments, 0))
		if err != nil {
			return b.errorValue(err)
		}
		result := object()
		result.Set("localize", b.function(func(_ js.Value, arguments []js.Value) any {
			return adapter.localize(arguments)
		}))
		result.Set("forRequest", b.function(func(_ js.Value, arguments []js.Value) any {
			return adapter.forRequest(valueArgument(arguments, 0))
		}))
		result.Set("languageFor", b.function(func(_ js.Value, arguments []js.Value) any {
			return adapter.languageFor(valueArgument(arguments, 0))
		}))
		return js.Global().Get("Object").Call("freeze", result)
	})
}

func (b *bridge) expressModule() js.Value {
	express := b.function(b.expressLocalizer)
	module := object()
	module.Set("default", express)
	module.Set("expressLocalizer", express)
	return module
}

func (b *bridge) expressLocalizer(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		options := valueArgument(arguments, 0)
		adapter, err := b.newRequestLocalizer(options)
		if err != nil {
			return b.errorValue(err)
		}
		autoLocalize := boolProperty(options, "autoLocalizeResponse")
		autoResponse := b.function(func(this js.Value, arguments []js.Value) any {
			request := valueArgument(arguments, 0)
			sendJSON := valueArgument(arguments, 1)
			localized := adapter.localizeWithRequest(arguments[2:], request)
			if isError(localized) {
				return localized
			}
			return sendJSON.Call("call", this, localized)
		})
		return b.function(func(_ js.Value, arguments []js.Value) any {
			request := valueArgument(arguments, 0)
			response := valueArgument(arguments, 1)
			next := valueArgument(arguments, 2)
			request.Set("localizeData", b.wrapped(adapter.forRequest(request)))

			if autoLocalize && response.Type() == js.TypeObject && response.Get("json").Type() == js.TypeFunction {
				sendJSON := response.Get("json")
				bound := autoResponse.Call("bind", response, request, sendJSON)
				response.Set("json", b.wrapped(bound))
			}
			if next.Type() == js.TypeFunction {
				next.Invoke()
			}
			return js.Undefined()
		})
	})
}

func (r *requestLocalizer) localizeWithRequest(arguments []js.Value, request js.Value) any {
	if len(arguments) == 0 {
		return r.bridge.typeError("data must be JSON-compatible")
	}
	header := r.languageFor(request)
	if len(arguments) > 1 && !isNullish(arguments[1]) {
		header = stringArgument(arguments, 1, "")
	}
	return r.bridge.localizeWith(r.engine, []js.Value{arguments[0], js.ValueOf(header)})
}

func (b *bridge) nestModule() js.Value {
	localizerConstructor := b.function(b.newNestDataLocalizer)
	interceptorConstructor := b.function(b.newNestInterceptor)
	createInterceptor := b.function(func(_ js.Value, arguments []js.Value) any {
		return b.newNestInterceptor(js.Undefined(), arguments)
	})
	module := object()
	module.Set("NestDataLocalizer", localizerConstructor)
	module.Set("DataLocalizerInterceptor", interceptorConstructor)
	module.Set("createNestInterceptor", createInterceptor)
	return module
}

func (b *bridge) newNestDataLocalizer(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		adapter, err := b.newRequestLocalizer(valueArgument(arguments, 0))
		if err != nil {
			return b.errorValue(err)
		}
		result := object()
		result.Set("localize", b.function(func(_ js.Value, arguments []js.Value) any {
			return adapter.localize(arguments)
		}))
		result.Set("forRequest", b.function(func(_ js.Value, arguments []js.Value) any {
			return adapter.forRequest(valueArgument(arguments, 0))
		}))
		return result
	})
}

func (b *bridge) newNestInterceptor(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		adapter, err := b.newRequestLocalizer(valueArgument(arguments, 0))
		if err != nil {
			return b.errorValue(err)
		}
		result := object()
		result.Set("intercept", b.function(func(_ js.Value, arguments []js.Value) any {
			return b.safely(func() any {
				context := valueArgument(arguments, 0)
				next := valueArgument(arguments, 1)
				request := context.Call("switchToHttp").Call("getRequest")
				localizeData := b.wrapped(adapter.forRequest(request))
				request.Set("localizeData", localizeData)
				require := js.Global().Get("__dataLocalizerRequire")
				mapOperator := require.Invoke("rxjs/operators").Get("map")
				return next.Call("handle").Call("pipe", mapOperator.Invoke(localizeData))
			})
		}))
		return result
	})
}

func (b *bridge) fastifyModule() js.Value {
	fastify := b.function(b.fastifyLocalizer)
	module := object()
	module.Set("default", fastify)
	module.Set("fastifyLocalizer", fastify)
	return module
}

func (b *bridge) fastifyLocalizer(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		fastify := valueArgument(arguments, 0)
		options := valueArgument(arguments, 1)
		done := valueArgument(arguments, 2)
		adapter, err := b.newRequestLocalizer(options)
		if err != nil {
			if done.Type() == js.TypeFunction {
				done.Invoke(b.errorValue(err))
				return js.Undefined()
			}
			return b.errorValue(err)
		}
		localizeData := b.wrapped(b.function(func(this js.Value, arguments []js.Value) any {
			return adapter.localizeWithRequest(arguments, this)
		}))
		fastify.Call("decorateRequest", "localizeData", localizeData)

		if boolProperty(options, "autoLocalizeResponse") {
			hook := b.wrapped(b.function(func(_ js.Value, arguments []js.Value) any {
				request := valueArgument(arguments, 0)
				payload := valueArgument(arguments, 2)
				next := valueArgument(arguments, 3)
				localized := adapter.localizeWithRequest([]js.Value{payload}, request)
				if isError(localized) {
					next.Invoke(localized)
				} else {
					next.Invoke(js.Null(), localized)
				}
				return js.Undefined()
			}))
			fastify.Call("addHook", "preSerialization", hook)
		}
		if done.Type() == js.TypeFunction {
			done.Invoke()
		}
		return js.Undefined()
	})
}

func (b *bridge) koaModule() js.Value {
	koa := b.function(b.koaLocalizer)
	module := object()
	module.Set("default", koa)
	module.Set("koaLocalizer", koa)
	return module
}

func (b *bridge) koaLocalizer(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		options := valueArgument(arguments, 0)
		adapter, err := b.newRequestLocalizer(options)
		if err != nil {
			return b.errorValue(err)
		}
		autoLocalize := boolProperty(options, "autoLocalizeResponse")
		afterNext := b.function(func(_ js.Value, arguments []js.Value) any {
			context := valueArgument(arguments, 0)
			request := valueArgument(arguments, 1)
			if autoLocalize && context.Get("body").Type() != js.TypeUndefined {
				localized := adapter.localizeWithRequest([]js.Value{context.Get("body")}, request)
				if isError(localized) {
					return localized
				}
				context.Set("body", localized)
			}
			return js.Undefined()
		})
		return b.function(func(_ js.Value, arguments []js.Value) any {
			return b.safely(func() any {
				context := valueArgument(arguments, 0)
				next := valueArgument(arguments, 1)
				request := context
				if candidate := context.Get("request"); !isNullish(candidate) {
					request = candidate
				}
				context.Set("localizeData", b.wrapped(adapter.forRequest(request)))
				after := b.wrapped(afterNext.Call("bind", js.Undefined(), context, request))
				result := next.Invoke()
				if !isNullish(result) && result.Get("then").Type() == js.TypeFunction {
					return result.Call("then", after)
				}
				return after.Invoke()
			})
		})
	})
}

func boolProperty(value js.Value, name string) bool {
	return !isNullish(value) && value.Type() == js.TypeObject && value.Get(name).Type() == js.TypeBoolean && value.Get(name).Bool()
}

func isError(value any) bool {
	jsValue, ok := value.(js.Value)
	return ok && jsValue.InstanceOf(js.Global().Get("Error"))
}
