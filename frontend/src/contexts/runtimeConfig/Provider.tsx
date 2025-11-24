import React, {
  createContext,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react'

import axiosInstance from '@/axiosConfig'

export type RuntimeConfig = {
  is_test_user_mode: boolean
  is_line_authentication: boolean
  version?: string
}

const DEFAULT_CONFIG: RuntimeConfig = {
  is_test_user_mode: false,
  is_line_authentication: false,
}

type StoredConfig = { ts: number; val: RuntimeConfig }

const STORAGE_KEY = 'runtime_cfg'
const CACHE_TTL_MS = 5 * 60 * 1000 // 5分

const loadFromStorage = (): RuntimeConfig | null => {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const parsed: StoredConfig = JSON.parse(raw)
    if (Date.now() - parsed.ts < CACHE_TTL_MS) return parsed.val
  } catch (e) {
    console.error('Failed to load runtime config from storage:', e)
  }
  return null
}

const saveToStorage = (config: RuntimeConfig): void => {
  try {
    const data: StoredConfig = { ts: Date.now(), val: config }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
  } catch (e) {
    console.error('Failed to save runtime config to storage:', e)
  }
}

const fetchRuntimeConfig = async (): Promise<RuntimeConfig> => {
  const { data } = await axiosInstance.get<RuntimeConfig>(
    '/public/runtime-config',
  )
  return data
}

export interface RuntimeConfigContextValue {
  config: RuntimeConfig
  isLoading: boolean
}

// eslint-disable-next-line react-refresh/only-export-components
export const RuntimeConfigContext = createContext<RuntimeConfigContextValue>({
  config: DEFAULT_CONFIG,
  isLoading: true,
})

export const RuntimeConfigProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [config, setConfig] = useState<RuntimeConfig>(DEFAULT_CONFIG)
  const [isLoading, setIsLoading] = useState(true)
  const lastFetchTimeRef = useRef<number>(0)

  const updateConfig = useCallback(
    (newConfig: RuntimeConfig, persist = true) => {
      setConfig(newConfig)
      if (persist) saveToStorage(newConfig)
      lastFetchTimeRef.current = Date.now()
    },
    [],
  )

  const refreshConfig = useCallback(async () => {
    try {
      const newConfig = await fetchRuntimeConfig()
      updateConfig(newConfig, true)
    } catch {
      const cached = loadFromStorage()
      if (cached) updateConfig(cached, false)
      else updateConfig(DEFAULT_CONFIG, false)
    }
  }, [updateConfig])

  // 初期化：localStorage 即時採用 → 背景で再取得
  useEffect(() => {
    const init = async () => {
      const cached = loadFromStorage()
      if (cached) {
        setConfig(cached)
        try {
          const raw = localStorage.getItem(STORAGE_KEY)
          if (raw) {
            const parsed: StoredConfig = JSON.parse(raw)
            lastFetchTimeRef.current = parsed.ts || Date.now()
          }
        } catch {
          /* noop */
        }
        setIsLoading(false)
        refreshConfig().catch(() => {})
      } else {
        await refreshConfig()
        setIsLoading(false)
      }
    }
    init()
  }, [refreshConfig])

  // 可視化されたら 5分超なら再取得
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.visibilityState !== 'visible') return
      const elapsed = Date.now() - lastFetchTimeRef.current
      if (elapsed >= CACHE_TTL_MS) {
        refreshConfig().catch(() => {})
      }
    }
    document.addEventListener('visibilitychange', handleVisibilityChange)
    return () =>
      document.removeEventListener('visibilitychange', handleVisibilityChange)
  }, [refreshConfig])

  return (
    <RuntimeConfigContext.Provider value={{ config, isLoading }}>
      {children}
    </RuntimeConfigContext.Provider>
  )
}
