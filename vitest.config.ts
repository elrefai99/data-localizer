import { defineConfig } from 'vitest/config'
import path from "node:path";

export default defineConfig({
  test: {
    env: {
      NODE_ENV: 'test'
    },
    globals: true,
    environment: 'node',
    include: ["src/**/*.{test,spec}.ts", "tests/**/*.{test,spec}.ts"],
    clearMocks: true,
    restoreMocks: true,
    coverage: {
      provider: 'v8',
      include: ['tests/**/*.ts'],
      exclude: ['tests/**/*.d.ts'],
    },
  },
  resolve: { alias: { "@": path.resolve(__dirname, "src") } },
})
