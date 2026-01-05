// eslint-disable-next-line simple-import-sort/imports
import { Page, expect } from '@playwright/test'

import { test as base } from './base'

export type UserRole = 'general' | 'admin' | 'root' | 'test'

type AuthProfile = {
  role: UserRole
  name: string
  token: string
}

type AuthFixtures = {
  signInAs: (
    role?: UserRole,
    options?: { persistToken?: boolean },
  ) => Promise<AuthProfile>
  goToMyPageFromSignIn: (role?: UserRole) => Promise<void>
  expectRoleProtectedNavigation: (role: UserRole) => Promise<void>
}

const authProfiles: Record<UserRole, AuthProfile> = {
  general: {
    role: 'general',
    name: 'E2E General User',
    token: 'e2e-token-general',
  },
  admin: {
    role: 'admin',
    name: 'E2E Admin User',
    token: 'e2e-token-admin',
  },
  root: {
    role: 'root',
    name: 'E2E Root User',
    token: 'e2e-token-root',
  },
  test: {
    role: 'test',
    name: 'E2E Test User',
    token: 'e2e-token-test',
  },
}

const APP_ORIGIN = 'http://localhost:5173'

const persistAuthToken = async (page: Page, token: string) => {
  await page.addInitScript(
    ([value]) => {
      localStorage.setItem('token', value)
      localStorage.removeItem('logoutMessage')
    },
    [token],
  )
  await page.evaluate(
    ([value]) => {
      localStorage.setItem('token', value)
      localStorage.removeItem('logoutMessage')
    },
    [token],
  )
}

const primeAuthMocks = async (page: Page, profile: AuthProfile) => {
  const { overrideAuthOnPage } = await import('../mocks/routeMocks')
  await overrideAuthOnPage(page, {
    status: 200,
    role: profile.role,
    userName: profile.name,
    token: profile.token,
  })
}

export const test = base.extend<AuthFixtures>({
  signInAs: async ({ page }, use) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(async (role = 'general', options = { persistToken: true }) => {
      const profile = authProfiles[role]
      await primeAuthMocks(page, profile)
      if (options.persistToken) {
        await persistAuthToken(page, profile.token)
      }
      return profile
    })
  },
  goToMyPageFromSignIn: async ({ page, signInAs }, use) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(async (role = 'general') => {
      const profile = await signInAs(role, { persistToken: false })
      await page.goto('/sign_in')
      await page.getByLabel('Email').fill('user@example.com')
      await page.getByLabel('Password').fill('password123')
      await page.getByRole('button', { name: 'メールでサインイン' }).click()

      await expect(page).toHaveURL(/\/mypage$/)
      await expect(
        page.getByRole('heading', { name: 'マイページ' }),
      ).toBeVisible()
      await expect(page.getByText(`${profile.name} さん`)).toBeVisible()
    })
  },
  expectRoleProtectedNavigation: async ({ page, signInAs }, use) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    await use(async (role: UserRole) => {
      await signInAs(role)
      const testCases = [
        {
          path: '/users',
          allowed: role === 'root',
          assertion: () =>
            expect(
              page.getByRole('heading', { name: 'ユーザー一覧' }),
            ).toBeVisible(),
        },
        {
          path: '/words/new',
          allowed: role === 'admin' || role === 'root',
          assertion: () =>
            expect(
              page.getByRole('heading', { name: '単語登録' }),
            ).toBeVisible(),
        },
      ] as const

      for (const testCase of testCases) {
        await page.goto(testCase.path)
        if (testCase.allowed) {
          await expect(page).toHaveURL(testCase.path)
          await testCase.assertion()
        } else {
          await expect(page).toHaveURL(`${APP_ORIGIN}/`)
        }
      }
    })
  },
})

export { expect }
