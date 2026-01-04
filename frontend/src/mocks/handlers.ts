import { rest, RestHandler, RestRequest } from 'msw'

const API_BASE_URL = 'http://localhost:8080'

export type AuthRole = 'general' | 'admin' | 'root' | 'test'
export type AuthStatus = 200 | 401 | 403

export type AuthMockOptions = {
  status?: AuthStatus
  role?: AuthRole
  userName?: string
  token?: string
}

export type MockOptions = {
  auth?: AuthMockOptions
}

type MockResponse = {
  status?: number
  body?: unknown
  headers?: Record<string, string>
}

export type MockDefinition = {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'OPTIONS'
  path: string
  resolver: (ctx: {
    params: Record<string, string>
    searchParams: URLSearchParams
    headers: Record<string, string | undefined>
    body?: unknown
    request: RestRequest
  }) => Promise<MockResponse> | MockResponse
}

const defaultUser = {
  id: 1,
  name: 'E2E User',
  isAdmin: false,
  isRoot: false,
  isTest: false,
}

const sampleWords = [
  {
    id: 1,
    name: 'apple',
    registrationCount: 5,
    wordInfos: [
      {
        id: 11,
        partOfSpeechId: 1,
        japaneseMeans: [{ id: 111, name: 'りんご' }],
      },
    ],
    isRegistered: false,
    attentionLevel: 1,
    quizCount: 3,
    correctCount: 2,
    memo: 'Red and sweet fruit',
  },
  {
    id: 2,
    name: 'banana',
    registrationCount: 4,
    wordInfos: [
      {
        id: 21,
        partOfSpeechId: 1,
        japaneseMeans: [{ id: 211, name: 'バナナ' }],
      },
    ],
    isRegistered: true,
    attentionLevel: 2,
    quizCount: 4,
    correctCount: 3,
    memo: 'Yellow and long fruit',
  },
]

const quizQuestions = [
  {
    quizID: 501,
    questionNumber: 1,
    wordName: 'apple',
    choicesJpms: [
      { japaneseMeanID: 1, name: 'りんご' },
      { japaneseMeanID: 2, name: 'バナナ' },
      { japaneseMeanID: 3, name: 'ぶどう' },
      { japaneseMeanID: 4, name: 'もも' },
    ],
  },
  {
    quizID: 501,
    questionNumber: 2,
    wordName: 'banana',
    choicesJpms: [
      { japaneseMeanID: 1, name: 'りんご' },
      { japaneseMeanID: 2, name: 'バナナ' },
      { japaneseMeanID: 3, name: 'ぶどう' },
      { japaneseMeanID: 4, name: 'もも' },
    ],
  },
]

const buildUserFromRole = (role: AuthRole | undefined) => {
  if (!role) return defaultUser
  return {
    ...defaultUser,
    isAdmin: role === 'admin' || role === 'root',
    isRoot: role === 'root',
    isTest: role === 'test',
  }
}

const buildResult = () => ({
  quizNumber: 501,
  totalQuestionsCount: quizQuestions.length,
  correctCount: quizQuestions.length,
  resultCorrectRate: 100,
  resultSetting: {
    quizSettingCompleted: true,
    questionCount: quizQuestions.length,
    isSaveResult: true,
    isRegisteredWords: 0,
    correctRate: 100,
    attentionLevelList: [1, 2, 3, 4, 5],
    partsOfSpeeches: [1, 3, 4, 5],
    isIdioms: 0,
    isSpecialCharacters: 0,
  },
  resultQuestions: quizQuestions.map((q, idx) => ({
    quizID: q.quizID,
    questionNumber: q.questionNumber,
    wordID: idx + 1,
    wordName: q.wordName,
    posID: 1,
    correctJpmID: q.choicesJpms[0].japaneseMeanID,
    choicesJpms: q.choicesJpms,
    answerJpmID: q.choicesJpms[0].japaneseMeanID,
    isCorrect: true,
    timeMs: 1500,
    registeredWord: {
      isRegistered: true,
      attentionLevel: 1,
      quizCount: 3,
      correctCount: 2,
    },
  })),
})

const buildAuthHandlers = (auth?: AuthMockOptions): MockDefinition[] => {
  const status = auth?.status ?? 200
  const role = auth?.role ?? 'general'
  const user = buildUserFromRole(role)
  const token = auth?.token ?? 'e2e-token'
  const bodyForStatus =
    status === 200
      ? { user }
      : {
          message: status === 401 ? 'ログインしてください' : '権限がありません',
        }

  return [
    {
      method: 'GET',
      path: '/auth/check',
      resolver: async ({ headers }) => {
        const authHeader = headers.authorization || headers.Authorization
        if (!authHeader && status === 200) {
          return { status: 401, body: { message: 'ログインしてください' } }
        }
        return { status, body: bodyForStatus }
      },
    },
    {
      method: 'GET',
      path: '/users/my_page',
      resolver: async ({ headers }) => {
        const authHeader = headers.authorization || headers.Authorization
        if (!authHeader && status === 200) {
          return { status: 401, body: { message: 'ログインしてください' } }
        }
        return {
          status,
          body:
            status === 200
              ? {
                  user: {
                    name: auth?.userName ?? user.name,
                    isAdmin: user.isAdmin,
                    isRoot: user.isRoot,
                    isTest: user.isTest,
                  },
                }
              : { message: bodyForStatus.message },
        }
      },
    },
    {
      method: 'OPTIONS',
      path: '/users/my_page',
      resolver: async () => ({ status: 200 }),
    },
    {
      method: 'POST',
      path: '/users/auth/test-login',
      resolver: async () => ({
        status: 200,
        body: {
          token,
          user_id: user.id,
          user_name: auth?.userName ?? user.name,
          jump: 'quiz',
        },
      }),
    },
    {
      method: 'POST',
      path: '/users/auth/test-logout',
      resolver: async () => ({ status: 200 }),
    },
    {
      method: 'POST',
      path: '/users/sign_in',
      resolver: async () => ({ status: 200, body: { token } }),
    },
    {
      method: 'POST',
      path: '/users/sign_up',
      resolver: async () => ({ status: 200, body: { token } }),
    },
  ]
}

export const buildMockDefinitions = (
  options?: MockOptions,
): MockDefinition[] => {
  const authHandlers = buildAuthHandlers(options?.auth)
  const runtimeHandlers: MockDefinition[] = [
    {
      method: 'GET',
      path: '/setting/user_config',
      resolver: async () => ({ status: 200, body: { is_dark_mode: false } }),
    },
    {
      method: 'GET',
      path: '/public/runtime-config',
      resolver: async () => ({
        status: 200,
        body: {
          is_test_user_mode: false,
          is_line_authentication: false,
          version: '1.0.0-e2e',
        },
      }),
    },
  ]

  const quizHandlers: MockDefinition[] = [
    {
      method: 'GET',
      path: '/quizzes',
      resolver: async () => ({
        status: 200,
        body: { isRunningQuiz: false, nextQuestion: quizQuestions[0] },
      }),
    },
    {
      method: 'POST',
      path: '/quizzes/new',
      resolver: async () => ({
        status: 200,
        body: {
          quizID: quizQuestions[0].quizID,
          totalCreateQuestion: quizQuestions.length,
          nextQuestion: quizQuestions[0],
        },
      }),
    },
    {
      method: 'POST',
      path: '/quizzes/answers/:quizID',
      resolver: async ({ body }) => {
        const payload = (body ?? {}) as { questionNumber?: number }
        const currentNumber = payload.questionNumber ?? 1
        const next = quizQuestions.find(
          (q) => q.questionNumber === currentNumber + 1,
        )
        if (next) {
          return {
            status: 200,
            body: {
              isFinish: false,
              isCorrect: true,
              nextQuestion: next,
              quizNumber: 501,
            },
          }
        }
        return {
          status: 200,
          body: {
            isFinish: true,
            isCorrect: true,
            quizNumber: 501,
          },
        }
      },
    },
    {
      method: 'GET',
      path: '/results',
      resolver: async () => ({
        status: 200,
        body: [
          {
            quizNumber: 501,
            createdAt: new Date().toISOString(),
            isRegisteredWords: 0,
            isIdioms: 0,
            isSpecialCharacters: 0,
            choicesPosIds: [1],
            totalQuestionsCount: quizQuestions.length,
            correctCount: quizQuestions.length,
            resultCorrectRate: 100,
          },
        ],
      }),
    },
    {
      method: 'GET',
      path: '/results/:quizNo',
      resolver: async () => ({ status: 200, body: buildResult() }),
    },
  ]

  const wordHandlers: MockDefinition[] = [
    {
      method: 'GET',
      path: '/words',
      resolver: async ({ searchParams }) => {
        const search = (searchParams.get('search') ?? '').toLowerCase()
        const filtered = search
          ? sampleWords.filter((w) => w.name.toLowerCase().includes(search))
          : sampleWords
        return { status: 200, body: { words: filtered, totalPages: 1 } }
      },
    },
    {
      method: 'GET',
      path: '/words/:id',
      resolver: async ({ params }) => {
        const id = Number(params.id)
        const word = sampleWords.find((w) => w.id === id)
        if (!word) return { status: 404, body: { message: 'Not found' } }
        return { status: 200, body: word }
      },
    },
    {
      method: 'POST',
      path: '/words/register',
      resolver: async ({ body }) => {
        const payload = (body ?? {}) as {
          wordId?: number
          IsRegistered?: boolean
        }
        const word = sampleWords.find((w) => w.id === payload.wordId)
        return {
          status: 200,
          body: {
            id: payload.wordId ?? 0,
            name: word?.name ?? 'word',
            isRegistered: payload.IsRegistered ?? false,
            registrationCount:
              (word?.registrationCount ?? 0) + (payload.IsRegistered ? 1 : -1),
          },
        }
      },
    },
    {
      method: 'POST',
      path: '/words/memo',
      resolver: async ({ body }) => {
        const payload = (body ?? {}) as { wordId?: number; memo?: string }
        return {
          status: 200,
          body: { wordId: payload.wordId, memo: payload.memo },
        }
      },
    },
    {
      method: 'DELETE',
      path: '/words/:id',
      resolver: async () => ({ status: 204 }),
    },
  ]

  return [...runtimeHandlers, ...authHandlers, ...quizHandlers, ...wordHandlers]
}

const createRestHandler = (def: MockDefinition): RestHandler => {
  const url = `${API_BASE_URL}${def.path}`
  const factory = (rest as Record<string, unknown>)[
    def.method.toLowerCase()
  ] as typeof rest.get

  if (typeof factory !== 'function') {
    throw new Error(`Unsupported HTTP method for MSW handler: ${def.method}`)
  }

  return factory(url, async (req, res, ctx) => {
    const headers = Object.fromEntries(req.headers.entries())
    const body = req.headers.get('content-type')?.includes('application/json')
      ? await req.json()
      : undefined
    const params = Object.fromEntries(
      Object.entries(req.params).map(([key, value]) => [
        key,
        Array.isArray(value) ? value[0] : value,
      ]),
    ) as Record<string, string>
    const result = await def.resolver({
      params,
      searchParams: new URL(req.url).searchParams,
      headers,
      body,
      request: req,
    })

    const status = result.status ?? 200
    const json = result.body ?? {}
    const headerEntries = Object.entries(result.headers ?? {}).map(([k, v]) =>
      ctx.set(k, v),
    )
    return res(ctx.status(status), ctx.json(json), ...headerEntries)
  })
}

export const createRestHandlers = (options?: MockOptions): RestHandler[] =>
  buildMockDefinitions(options).map(createRestHandler)

export const createAuthHandlers = (auth?: AuthMockOptions): RestHandler[] =>
  buildAuthHandlers(auth).map(createRestHandler)

export const buildAuthDefinitions = (
  auth?: AuthMockOptions,
): MockDefinition[] => buildAuthHandlers(auth)

export { API_BASE_URL }
