// src/pages/legal/LegalLayout.tsx
import * as React from 'react'

import { PageContainer } from '@/components/ui/card' // パス調整
import TextArticle from '@/components/ui/TextArticle'

type Props = { title: string; updated: string; children: React.ReactNode }

export default function LegalLayout({ title, updated, children }: Props) {
  return (
    <PageContainer>
      <TextArticle title={title} updated={updated}>
        {children}
      </TextArticle>
    </PageContainer>
  )
}
