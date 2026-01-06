// eslint-disable-next-line simple-import-sort/imports
import { test as base, expect } from '@playwright/test'

import {
  applyDefaultApiMocks,
  AuthScenario,
  overrideAuthOnPage,
} from '../mocks/routeMocks.ts'

type Fixtures = {
  useAuthMock: (scenario?: AuthScenario) => Promise<void>
}

export const test = base.extend<Fixtures>({
  page: async ({ page }, use) => {
    await applyDefaultApiMocks(page)
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(page)
  },
  useAuthMock: async ({ page }, use) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(async (scenario = 'authorized') => {
      await overrideAuthOnPage(page, scenario)
    })
  },
})

export { expect }
