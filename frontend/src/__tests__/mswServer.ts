import { setupServer } from 'msw/node'

import { createRestHandlers } from '../../src/mocks/handlers'

export const server = setupServer(...createRestHandlers())
