import path from 'path'

import react from '@vitejs/plugin-react-swc'
import { defineConfig, loadEnv } from 'vite'
import checker from 'vite-plugin-checker'

export default defineConfig(({ mode }) => {
  /* .env / Vercel の環境変数を読み込む */
  const env = loadEnv(mode, process.cwd()) // mode = 'development' | 'production'

  return {
    base: '/', // ← Vercel 配信はこれでOK（サブパス配信ならそのパスに）
    build: {
      outDir: 'dist', // ← Vercel の dist と一致
      emptyOutDir: true,
      assetsDir: 'assets', // 既定だが明示しておくと安心
    },
    plugins: [
      react(),
      checker({
        typescript: {
          buildMode: false, // ビルド時にはTypeScriptチェックを無効化
        },
        eslint: {
          // Flat Config を使うことを明示
          useFlatConfig: true,
          // ← ここ、元設定は ' "..."' になっていて eslint コマンド名が抜けていたので修正
          lintCommand: 'eslint "./src/**/*.{js,jsx,ts,tsx}"',
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
        '@': path.resolve(__dirname, 'src'),
      },
    },
    /* ビルド時に API URL を注入 */
    define: {
      __API_URL__: JSON.stringify(env.VITE_API_URL ?? ''),
    },
  }
})
