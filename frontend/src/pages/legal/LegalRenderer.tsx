// src/pages/legal/LegalRenderer.tsx
import * as React from 'react'

import type { Node } from './legalContent'

export function LegalRenderer({ nodes }: { nodes: Node[] }) {
  return (
    <>
      {nodes.map((node, i) => {
        switch (node.type) {
          case 'h1':
            return <h1 key={i}>{node.text}</h1>
          case 'h2':
            return <h2 key={i}>{node.text}</h2>
          case 'p':
            return <p key={i}>{node.text}</p>
          case 'ul':
            return (
              <ul key={i}>
                {node.items.map((it, j) => (
                  <li key={j}>{it}</li>
                ))}
              </ul>
            )
          case 'link':
            return (
              <p key={i}>
                <a href={node.href} target="_blank" rel="noreferrer">
                  {node.text}
                </a>
              </p>
            )
          case 'blockquote':
            return <blockquote key={i}>{node.text}</blockquote>
          case 'hr':
            return <hr key={i} />
          default:
            return null
        }
      })}
    </>
  )
}
