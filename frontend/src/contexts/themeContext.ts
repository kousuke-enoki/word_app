import { createContext, useContext } from 'react'

export type Theme = 'light' | 'dark'
export interface ThemeCtx {
  theme: Theme
  setTheme: (t: Theme) => void
}

export const ThemeContext = createContext<ThemeCtx>({
  theme: 'light',
  setTheme: () => {},
})

export const useTheme = () => useContext(ThemeContext)
