import { Page, Route } from '@playwright/test'

import {
  API_BASE_URL,
  AuthMockOptions,
  buildAuthDefinitions,
  buildMockDefinitions,
  MockDefinition,
  MockOptions,
} from '../../../src/mocks/handlers'

type RouteRegistration = {
  matcher: RegExp
  handler: (route: Route) => Promise<void>
}

export type AuthScenario = 'authorized' | 'unauthorized' | 'forbidden'

const scenarioToAuthOptions: Record<AuthScenario, AuthMockOptions> = {
  authorized: { status: 200 },
  unauthorized: { status: 401 },
  forbidden: { status: 403 },
}

const escapeRegex = (value: string) =>
  value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

const buildMatcher = (path: string) => {
  const paramNames: string[] = []
  const escaped = escapeRegex(path).replace(
    /\\:([A-Za-z0-9_]+)/g,
    (_m, key) => {
      paramNames.push(key)
      return '([^/]+)'
    },
  )
  const matcher = new RegExp(`^${escapeRegex(API_BASE_URL)}${escaped}$`)
  return { matcher, paramNames }
}

const extractParams = (
  matcher: RegExp,
  url: string,
  paramNames: string[],
): Record<string, string> => {
  const match = matcher.exec(url)
  if (!match) return {}
  const [, ...groups] = match
  return paramNames.reduce<Record<string, string>>((acc, name, idx) => {
    acc[name] = groups[idx]
    return acc
  }, {})
}

const parseBody = (route: Route) => {
  const postData = route.request().postData()
  if (!postData) return undefined
  try {
    return JSON.parse(postData)
  } catch {
    return postData
  }
}

const registerRoutes = async (
  page: Page,
  definitions: MockDefinition[],
): Promise<RouteRegistration[]> => {
  const registrations: RouteRegistration[] = []
  for (const def of definitions) {
    const { matcher, paramNames } = buildMatcher(def.path)
    const handler = async (route: Route) => {
      if (route.request().method() !== def.method) return route.fallback()
      const url = new URL(route.request().url())
      const params = extractParams(matcher, url.href, paramNames)
      const result = await def.resolver({
        params,
        searchParams: url.searchParams,
        headers: route.request().headers(),
        body: parseBody(route),
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        request: route.request() as any,
      })
      const headers = {
        'content-type': 'application/json',
        ...(result.headers ?? {}),
      }

      if (result.body === undefined) {
        await route.fulfill({ status: result.status ?? 200, headers })
        return
      }

      await route.fulfill({
        status: result.status ?? 200,
        body: JSON.stringify(result.body),
        headers,
      })
    }
    await page.route(matcher, handler)
    registrations.push({ matcher, handler })
  }
  return registrations
}

export const applyDefaultApiMocks = async (
  page: Page,
  options?: MockOptions,
) => {
  await registerRoutes(page, buildMockDefinitions(options))
}

export const overrideAuthOnPage = async (
  page: Page,
  scenario: AuthScenario | AuthMockOptions,
) => {
  const opts =
    typeof scenario === 'string' ? scenarioToAuthOptions[scenario] : scenario
  await registerRoutes(page, buildAuthDefinitions(opts))
}
