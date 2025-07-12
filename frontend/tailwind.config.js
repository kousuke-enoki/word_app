/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './index.html',                  // Vite のエントリ
    './src/**/*.{js,ts,jsx,tsx}',    // すべての React/TS ファイル
  ],
  theme: {
    extend:{
      colors:{ primary:'#2563eb', surface:'#f9fafb' },
      borderRadius:{ card:'0.75rem' }
    }
  },
  plugins: [],
};
