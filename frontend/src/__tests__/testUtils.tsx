import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { QueryClientConfig } from '@tanstack/react-query'
import { render, type RenderOptions } from '@testing-library/react'
import React from 'react'
import type { i18n as I18nInstance } from 'i18next'
import { I18nextProvider } from 'react-i18next'
import {
  MemoryRouter,
  type MemoryRouterProps,
  Route,
  Routes,
} from 'react-router-dom'
import { afterEach } from 'vitest'

type RouterOptions = {
  /** MemoryRouter の initialEntries */
  initialEntries?: MemoryRouterProps['initialEntries']
  /** MemoryRouter の initialIndex */
  initialIndex?: MemoryRouterProps['initialIndex']
  /** <Routes> に渡すルート定義。省略時は ui を "*" で描画 */
  routes?: Array<{ path: string; element: React.ReactElement }>
}

type RenderWithClientOptions = {
  /** 必要ならカスタム QueryClient を差し込む */
  queryClient?: QueryClient
  /** i18n インスタンスを渡すと I18nextProvider で wrap */
  i18n?: I18nInstance
  /** MemoryRouter の設定 */
  router?: RouterOptions
  /** さらに独自 Provider で wrap したい場合 */
  wrapper?: React.ComponentType<{ children: React.ReactNode }>
  /** testing-library の render オプション */
  renderOptions?: Omit<RenderOptions, 'wrapper'>
  /** QueryClient のデフォルト設定上書き */
  queryClientConfig?: QueryClientConfig
}

/** テストごとに QueryClient を用意し、afterEach で clear() する */
export const createTestQueryClient = (config?: QueryClientConfig) => {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false } },
    ...config,
  })

  afterEach(() => client.clear())

  return client
}

/** Provider 付き render */
export const renderWithClient = (
  ui: React.ReactElement,
  options: RenderWithClientOptions = {},
) => {
  const client = options.queryClient ?? createTestQueryClient(options.queryClientConfig)

  const routes = options.router?.routes ?? [{ path: '*', element: ui }]

  const routedUi = options.router ? (
    <MemoryRouter
      initialEntries={options.router.initialEntries}
      initialIndex={options.router.initialIndex}
    >
      <Routes>
        {routes.map((route) => (
          <Route key={route.path} path={route.path} element={route.element} />
        ))}
      </Routes>
    </MemoryRouter>
  ) : (
    ui
  )

  const withI18n = options.i18n ? (
    <I18nextProvider i18n={options.i18n}>{routedUi}</I18nextProvider>
  ) : (
    routedUi
  )

  const wrappedUi = options.wrapper ? (
    <options.wrapper>{withI18n}</options.wrapper>
  ) : (
    withI18n
  )

  return {
    client,
    ...render(
      <QueryClientProvider client={client}>{wrappedUi}</QueryClientProvider>,
      options.renderOptions,
    ),
  }
}
