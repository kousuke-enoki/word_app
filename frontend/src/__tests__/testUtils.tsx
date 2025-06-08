import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render } from '@testing-library/react'
import React from 'react'

/** ★ 共有用に export する */
export let queryClient: QueryClient

/** Provider 付き render */
export const renderWithClient = (ui: React.ReactElement) => {
  // テストごとに新しいインスタンスを作成
  queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })

  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>,
  )
}