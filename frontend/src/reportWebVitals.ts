import type { ReportHandler } from 'web-vitals'    // 型だけなら OK

export const reportWebVitals = (onPerfEntry?: ReportHandler) => {
  if (!onPerfEntry) return

  import('web-vitals').then(({ onCLS, onFCP, onLCP, onTTFB, onINP }) => {
    onCLS(onPerfEntry)
    onFCP(onPerfEntry)
    onLCP(onPerfEntry)
    onTTFB(onPerfEntry)
    onINP?.(onPerfEntry)
  })
}
