# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type aware lint rules:

- Configure the top-level `parserOptions` property like this:

```js
export default tseslint.config({
  languageOptions: {
    // other options...
    parserOptions: {
      project: ['./tsconfig.node.json', './tsconfig.app.json'],
      tsconfigRootDir: import.meta.dirname,
    },
  },
})
```

- Replace `tseslint.configs.recommended` to `tseslint.configs.recommendedTypeChecked` or `tseslint.configs.strictTypeChecked`
- Optionally add `...tseslint.configs.stylisticTypeChecked`
- Install [eslint-plugin-react](https://github.com/jsx-eslint/eslint-plugin-react) and update the config:

```js
// eslint.config.js
import react from 'eslint-plugin-react'

export default tseslint.config({
  // Set the react version
  settings: { react: { version: '18.3' } },
  plugins: {
    // Add the react plugin
    react,
  },
  rules: {
    // other rules...
    // Enable its recommended rules
    ...react.configs.recommended.rules,
    ...react.configs['jsx-runtime'].rules,
  },
})
```

## Integration tests

- Run all integration suites with `pnpm vitest run "**/*.integration.test.tsx"`. To focus on one scenario, point to a specific file such as `pnpm vitest run src/routes/__tests__/auth.integration.test.tsx`.
- Name integration specs with the `*.integration.test.tsx` suffix. Place them close to the feature they cover (ideally next to the component/page).
- Tests share the MSW server defined in `src/__tests__/mswServer.ts`. Override handlers inside a test with `server.use(...)` as needed, for example:

```ts
import { rest } from 'msw'
import { server } from '../mswServer'

server.use(
  rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
    res(ctx.status(500)),
  ),
)
```

`server.resetHandlers()` runs after each test (via `src/__tests__/setupTests.ts`), so reapply overrides in every test case that needs them.
