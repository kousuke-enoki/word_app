// src/features/settings/useTestUserMode.ts
import { useEffect, useRef, useState } from 'react'

type RootCfg = { is_test_user_mode: boolean }
const CACHE_KEY = 'root_cfg_cache'
const CACHE_TTL_MS = 60_000 // 60秒ごとに再取得

function readCache(): { data: RootCfg; at: number } | null {
  try {
    const raw = sessionStorage.getItem(CACHE_KEY)
    if (!raw) return null
    return JSON.parse(raw)
  } catch {
    return null
  }
}

// function writeCache(data: RootCfg) {
//   sessionStorage.setItem(CACHE_KEY, JSON.stringify({ data, at: Date.now() }))
// }

export function useTestUserMode() {
  const [enabled, setEnabled] = useState<boolean>(() => {
    const c = readCache()
    return c?.data?.is_test_user_mode ?? false
  })
  const timerRef = useRef<number | null>(null)

  useEffect(() => {
    let mounted = true
    const fetchCfg = async () => {
      try {
        // キャッシュ新鮮ならスキップ
        // const cached = readCache()
        // if (cached && Date.now() - cached.at < CACHE_TTL_MS) {
        //   if (mounted) setEnabled(cached.data.is_test_user_mode)
        //   return
        // }
        // const { data } = await axiosInstance.get<RootCfg>('/settings/root')
        // writeCache(data)
        // if (mounted) setEnabled(!!data.is_test_user_mode)
        if (mounted) setEnabled(true) //テスト用でいったん常にtrue
      } catch {
        // 失敗時はキャッシュ/現状のまま
      }
    }

    // 初回＆定期ポーリング
    fetchCfg()
    timerRef.current = window.setInterval(fetchCfg, CACHE_TTL_MS)

    return () => {
      mounted = false
      if (timerRef.current) window.clearInterval(timerRef.current)
    }
  }, [])

  return enabled
}
