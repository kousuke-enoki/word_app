// frontend/src/types/wordTypes.ts

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
  partOfSpeech: PartOfSpeech
  japaneseMeans: JapaneseMean[]
}

export interface Word {
  id: number
  name: string
  wordInfos: WordInfo[]
  isRegistered?: boolean
  testCount?: number
  checkCount?: number
  registrationActive?: boolean
  memo?: string
}
