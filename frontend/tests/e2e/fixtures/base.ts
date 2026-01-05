// eslint-disable-next-line simple-import-sort/imports
import { test as base, expect } from '@playwright/test'

// AuthScenario 型を直接定義（routeMocks.ts からの型インポートを削除）
type AuthScenario = 'authorized' | 'unauthorized' | 'forbidden'

type Fixtures = {
  useAuthMock: (scenario?: AuthScenario) => Promise<void>
}

export const test = base.extend<Fixtures>({
  page: async ({ page }, use) => {
    const { applyDefaultApiMocks } = await import('../mocks/routeMocks.ts')
    await applyDefaultApiMocks(page)
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(page)
  },
  useAuthMock: async ({ page }, use) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(async (scenario = 'authorized') => {
      const { overrideAuthOnPage } = await import('../mocks/routeMocks.ts')
      await overrideAuthOnPage(page, scenario)
    })
  },
})

export { expect }
