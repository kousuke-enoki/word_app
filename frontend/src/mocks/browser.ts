import { setupWorker } from 'msw'

import { createRestHandlers } from './handlers'

export const worker = setupWorker(...createRestHandlers())
