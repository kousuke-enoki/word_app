// src/features/settings/useTestUserMode.ts
import { useRuntimeConfig } from '@/contexts/runtimeConfig/useRuntimeConfig'

/**
 * RuntimeConfig から is_test_user_mode を返す薄いラッパー。
 * 既存コードの置換用（互換API）。
 */
export function useTestUserMode(): boolean {
  const { config } = useRuntimeConfig()
  return !!config.is_test_user_mode
}
