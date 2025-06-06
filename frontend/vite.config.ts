import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: true,
    hmr: {
      protocol: 'wss',
      host: 'grouper-quality-vervet.ngrok-free.app',
      clientPort: 443,       // ← ブラウザ用。Vite 自体は 3000 で Listen
    },
    proxy: {
      // フロントからは /api/xxx で呼ぶ
      '/api': {
        target: 'http://backend:8080', // ← docker-compose の service 名
        changeOrigin: true,            // ← Host ヘッダを backend:8080 に書き換え
        secure: false,                 // (HTTPS を介さないので false)
        // rewrite: (p) => p.replace(/^\//, ''),
      },
    },
  },
  resolve: {
    alias: {
      react: path.resolve(__dirname, 'node_modules/react'),
      'react-dom': path.resolve(__dirname, 'node_modules/react-dom'),
      '@': path.resolve(__dirname, 'src'),
    },
  },
});
