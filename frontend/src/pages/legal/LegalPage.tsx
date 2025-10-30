// src/pages/legal/LegalPage.tsx
import * as React from 'react'

import { legalContent } from './legalContent'
import LegalLayout from './LegalLayout'
import { LegalRenderer } from './LegalRenderer'

type Slug = keyof typeof legalContent

export default function LegalPage({ slug }: { slug: Slug }) {
  const doc = legalContent[slug]
  return (
    <LegalLayout title={doc.title} updated={doc.updated}>
      <LegalRenderer nodes={doc.nodes} />
    </LegalLayout>
  )
}
