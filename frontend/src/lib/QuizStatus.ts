export const QuizStatus = {
  registered: { 0: '全単語', 1: 'のみ' } as const,
  idioms:     { 0: 'なし',   1: '含む' } as const,
  special:    { 0: 'なし',   1: '含む' } as const,
} as const;

export type QuizStatusKey = keyof typeof QuizStatus;

export const labelOf = <K extends QuizStatusKey>(
  k: K,
  v: keyof typeof QuizStatus[K],
) => QuizStatus[k][v];
