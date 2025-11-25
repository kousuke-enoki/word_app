// src/components/ui/TextArticle.tsx
import * as React from 'react'

import { Card } from './card' // ← Card を定義しているファイルに合わせてパス調整

type Props = {
  title: string
  updated?: string
  className?: string
  children: React.ReactNode
}

export default function TextArticle({
  title,
  updated,
  className = '',
  children,
}: Props) {
  return (
    <Card className={`p-6 md:p-10 ${className}`}>
      {/* prose は内部に独自の色/余白を持つので、dark では反転 */}
      <article
        className="
          prose prose-slate dark:prose-invert max-w-none text-left
          /* 見出しの上下を全ページで統一 */
          prose-headings:mt-8 prose-headings:mb-3

          /* p / ul / ol の“本文ブロック”としての上下余白＆行間を統一 */
          prose-p:my-4 prose-p:leading-7
          prose-ul:my-4
          prose-ol:my-4
          prose-li:leading-7        /* li も p と同じ行間 */

          /* リンク */
          prose-a:no-underline hover:prose-a:underline
          prose-a:text-blue-600 dark:prose-a:text-blue-400
        "
      >
        {/* not-prose でヘッダ部だけ自由にレイアウト */}
        <header className="not-prose mb-6">
          <h1 className="text-2xl md:text-3xl font-semibold tracking-tight">
            {title}
          </h1>
          {updated && (
            <p className="mt-1 text-xs opacity-70">最終更新日：{updated}</p>
          )}
        </header>
        {children}
      </article>
    </Card>
  )
}
