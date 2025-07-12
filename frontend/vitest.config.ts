// vitest.config.ts
import { defineConfig } from 'vitest/config'
import path from 'node:path'

export default defineConfig({
  test: {
    environment: 'jsdom',          // ← これで document/Window を注入
    globals: true,                 // beforeAll/vi 等をグローバルで使う
    setupFiles: ['./src/__tests__/setupTests.ts'],  // ← 1 行で閉じる！パスも正しく
    alias: {                       // Vite の @ を共有
      '@': path.resolve(__dirname, './src'),
    },
  },
})
