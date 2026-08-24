package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "buildnpm:", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := findModuleRoot()
	if err != nil {
		return err
	}
	destination := filepath.Join(root, "dist")
	if err := os.RemoveAll(destination); err != nil {
		return err
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}

	wasmPath := filepath.Join(destination, "data-localizer.wasm")
	command := exec.Command("go", "build", "-trimpath", "-o", wasmPath, "./cmd/wasm")
	command.Dir = root
	command.Env = buildEnvironment(os.Environ(), map[string]string{
		"CGO_ENABLED": "0",
		"GOARCH":      "wasm",
		"GOOS":        "js",
	})
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("compile WebAssembly: %w", err)
	}

	wasmExec, err := findWasmExec()
	if err != nil {
		return err
	}
	files := [][2]string{
		{filepath.Join(root, "bindings", "node", "index.js"), filepath.Join(destination, "index.js")},
		{filepath.Join(root, "bindings", "node", "index.d.ts"), filepath.Join(destination, "index.d.ts")},
		{wasmExec, filepath.Join(destination, "wasm_exec.js")},
		{filepath.Join(runtime.GOROOT(), "LICENSE"), filepath.Join(destination, "GO-LICENSE")},
	}
	for _, name := range []string{
		"shared.js",
		"framework.js", "framework.d.ts",
		"express.js", "express.d.ts",
		"nest.js", "nest.d.ts",
		"fastify.js", "fastify.d.ts",
		"koa.js", "koa.d.ts",
	} {
		files = append(files, [2]string{
			filepath.Join(root, "bindings", "node", "adapters", name),
			filepath.Join(destination, "adapters", name),
		})
	}
	for _, file := range files {
		if err := copyFile(file[0], file[1]); err != nil {
			return err
		}
	}
	return nil
}

func findModuleRoot() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
			return directory, nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("could not find go.mod")
		}
		directory = parent
	}
}

func findWasmExec() (string, error) {
	candidates := []string{
		filepath.Join(runtime.GOROOT(), "lib", "wasm", "wasm_exec.js"),
		filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("wasm_exec.js was not found under GOROOT %s", runtime.GOROOT())
}

func buildEnvironment(current []string, overrides map[string]string) []string {
	result := make([]string, 0, len(current)+len(overrides))
	for _, entry := range current {
		name, _, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		if _, replaced := overrides[strings.ToUpper(name)]; !replaced {
			result = append(result, entry)
		}
	}
	for name, value := range overrides {
		result = append(result, name+"="+value)
	}
	return result
}

func copyFile(source, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		return err
	}
	return output.Close()
}
