export const QuizStatus = {
  registered: { 0: '全単語', 1: 'のみ' } as const,
  idioms:     { 0: '全て',   1: '含む', 2: '含まない'} as const,
  special:    { 0: 'なし',   1: '含む', 2: '含まない'} as const,
} as const;

export type QuizStatusKey = keyof typeof QuizStatus;

export const labelOf = <K extends QuizStatusKey>(
  k: K,
  v: keyof typeof QuizStatus[K],
) => QuizStatus[k][v];
