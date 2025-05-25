export interface PartOfSpeech {
  id: number
  name: string
}

export interface JapaneseMean {
  id: number
  name: string
}

export interface WordInfo {
  id: number
  // partOfSpeech: PartOfSpeech
  partOfSpeechId: number
  japaneseMeans: JapaneseMean[]
}

export interface Word {
  id: number
  name: string
  registrationCount: number
  wordInfos: WordInfo[]
  isRegistered: boolean
  attentionLevel: number
  QuizCount?: number
  CorrectCount?: number
  memo?: string
}
