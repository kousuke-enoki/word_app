import { defineConfig } from 'vitest/config';   // ← そのまま
import { loadEnv } from 'vite';                 // ★ 追加
import react from '@vitejs/plugin-react-swc';
import path from 'path';
import checker from 'vite-plugin-checker';

export default defineConfig(({ mode }) => {
  /* .env / Vercel の環境変数を読み込む */
  const env = loadEnv(mode, process.cwd());      // mode = 'development' | 'production'

  return {
    plugins: [
      react(),
      checker({
        typescript: true,
        eslint: {
          lintCommand: ' "./src/**/*.{js,jsx,ts,tsx}"',
        },
      }),
    ],
    /* dev サーバ設定 (Vercel では無視される) */
    server: {
      port: 3000,
      host: '0.0.0.0',
    },
    resolve: {
      alias: {
        react: path.resolve(__dirname, 'node_modules/react'),
        'react-dom': path.resolve(__dirname, 'node_modules/react-dom'),
        '@': path.resolve(__dirname, 'src'),
      },
    },
    /* ★ ビルド時に API URL を注入 */
    define: {
      __API_URL__: JSON.stringify(env.VITE_API_URL ?? ''),
    },
    /* Vitest 設定 */
    test: {
      environment: 'jsdom',
      globals: true,
      setupFiles: './src/__tests__/setupTests.ts',
      css: true,
      coverage: {
        provider: 'v8',
        reporter: ['text', 'lcov'],
        reportsDirectory: './coverage'
      }
    },
  };
});
