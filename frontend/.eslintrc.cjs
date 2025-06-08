module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: { ecmaVersion: 'latest', sourceType: 'module' },
  plugins: ['@typescript-eslint', 'react-hooks'],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
    // React 19 の new JSX runtime に合わせて
    'plugin:react/recommended'
  ],
  settings: {
    react: { version: 'detect' },
  },
  env: { browser: true, es2023: true, node: true },
  ignorePatterns: ['dist', 'coverage'],
};
