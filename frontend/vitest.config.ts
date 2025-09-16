// vitest.config.ts
import path from 'node:path'

import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    environment: 'jsdom', // ← これで document/Window を注入
    globals: true, // beforeAll/vi 等をグローバルで使う
    setupFiles: ['./src/__tests__/setupTests.ts'], // ← 1 行で閉じる！パスも正しく
  },
  resolve: {
    alias: {
      // Vite の @ を共有
      '@': path.resolve(__dirname, './src'),
    },
  },
})
