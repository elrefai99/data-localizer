package main

import (
	"strings"
	"testing"
)

func TestGeneratedFiles(t *testing.T) {
	want := []string{
		"runtime.js",
		"index.js",
		"index.d.ts",
		"adapters/framework.js",
		"adapters/framework.d.ts",
		"adapters/express.js",
		"adapters/express.d.ts",
		"adapters/nest.js",
		"adapters/nest.d.ts",
		"adapters/fastify.js",
		"adapters/fastify.d.ts",
		"adapters/koa.js",
		"adapters/koa.d.ts",
	}

	files := generatedFiles()
	if len(files) != len(want) {
		t.Fatalf("generatedFiles() returned %d files, want %d", len(files), len(want))
	}
	for _, name := range want {
		contents, ok := files[name]
		if !ok {
			t.Errorf("generatedFiles() is missing %q", name)
			continue
		}
		if !strings.HasPrefix(contents, generatedNotice) {
			t.Errorf("generated file %q has no generated notice", name)
		}
		if strings.Contains(contents, "bindings/") || strings.Contains(contents, `bindings\`) {
			t.Errorf("generated file %q depends on bindings", name)
		}
	}
}

func TestRuntimeContainsOnlyBridgeCode(t *testing.T) {
	runtime := generateRuntime()
	for _, implementation := range []string{
		"createRequestLocalizer",
		"expressLocalizer",
		"fastifyLocalizer",
		"koaLocalizer",
		"getAcceptLanguage",
	} {
		if strings.Contains(runtime, implementation) {
			t.Errorf("generated runtime contains application implementation %q", implementation)
		}
	}
}
