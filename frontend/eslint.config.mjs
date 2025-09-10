import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import unusedImports from 'eslint-plugin-unused-imports'
import simpleImportSort from 'eslint-plugin-simple-import-sort'
import prettier from 'eslint-config-prettier'

// Flat Config では「上から順にマージ」される。
// 旧 .eslintrc の "extends" 相当は、ここで配列として“並べる”こと。
export default tseslint.config(
  { ignores: ['dist'] },
  // 旧: extends: [js.configs.recommended, ...tseslint.configs.recommended, prettier]
  js.configs.recommended,
  ...tseslint.configs.recommended,
  prettier,
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2020,
      // TSパーサを明示（tseslint.configs.recommendedでも入るが明示しておくと安心）
      parser: tseslint.parser,
      parserOptions: { sourceType: 'module' },
      globals: globals.browser,
    },
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
      'unused-imports': unusedImports,
      'simple-import-sort': simpleImportSort,
    },
    rules: {
      // React hooks
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],

      // 未使用 import/変数
      'unused-imports/no-unused-imports': 'error',
      'unused-imports/no-unused-vars': [
        'warn',
        { vars: 'all', args: 'after-used', ignoreRestSiblings: true },
      ],

      // import順（任意）
      'simple-import-sort/imports': 'error',
      'simple-import-sort/exports': 'error',
    },
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      prettier,
    ],
  },
)
