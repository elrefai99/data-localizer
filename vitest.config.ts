import { defineConfig } from 'vitest/config'
import path from "node:path";

export default defineConfig({
  test: {
    env: {
      NODE_ENV: 'test'
    },
    globals: true,
    environment: 'node',
    include: ["src/**/*.{test,spec}.ts"],
    clearMocks: true,
    restoreMocks: true,
    coverage: {
      provider: 'v8',
      include: ['src/**/*.ts'],
      exclude: ['src/generated/**', 'src/**/*.d.ts'],
    },
  },
  resolve: { alias: { "@": path.resolve(__dirname, "src") } },
})
