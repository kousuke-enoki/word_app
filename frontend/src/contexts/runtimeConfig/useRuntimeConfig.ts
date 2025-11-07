import { useContext } from 'react'

import {
  RuntimeConfigContext,
  type RuntimeConfigContextValue,
} from './Provider'

export const useRuntimeConfig = (): RuntimeConfigContextValue => {
  return useContext(RuntimeConfigContext)
}
