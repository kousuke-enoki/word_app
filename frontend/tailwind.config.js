/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './index.html',                  // Vite のエントリ
    './src/**/*.{js,ts,jsx,tsx}',    // すべての React/TS ファイル
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};
