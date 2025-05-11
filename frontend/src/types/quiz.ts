// src/types/quiz.ts
export interface CreateQuizReq {
  questionCount: number;
  isSaveResult: boolean;
  isRegisteredWords: number;
  correctRate: number;
  attentionLevelList: number[];
  partsOrSpeeches: number[];
  isIdioms: number;
  isSpecialCharacters: number;
}

export interface CreateQuizResponse {
  quizID: number;
  totalCreateQuestion: number;
  nextQuestion: QuizQuestion;
}

export interface PostAnswerQuestionRequest {
  quizID: number;
  questionNumber: number;
  answerJpmID: number;
}

export interface GetQuizRequest {
  quizID: number;
  beforeQuestionNumber: number;
}

export interface GetQuizResponce {
  isRunningQuiz: boolean;
  nextQuestion: QuizQuestion;
}

export interface QuizQuestion {
  quizID: number;
  questionNumber: number;
  wordName: string;
  choicesJpms: ChoiceJpm[];
}

export interface ChoiceJpm {
  japaneseMeanID: number;
  name: string;
}

/** POST /quizzes/:id/answers のレスポンス想定 */
export interface AnswerRouteRes {
  isFinish: boolean;
  isCorrect: boolean;
  nextQuestion?: QuizQuestion;     // isFinish === false
  result?: ResultRes;                    // isFinish === true
}

export interface ResultRes {
  quizNumber: number;
  totalQuestionsCount: number;
  correctCount: number;
  resultCorrectRate: number;
  resultSetting: QuizSettingsType;
  resultQuestions: ResultQuestion[];
}

export interface ResultQuestion {
  quizID: number;
  questionNumber: number;
  wordID: number;
  wordName: string;
  posID: number;
  correctJpmId: number;
  choicesJpms: ChoiceJpm[];
  answerJpmId: number;
  isCorrect: boolean;
  timeMs: number;
  registeredWord: registeredWord;
}

export interface registeredWord {
  isRegistered: boolean
  attentionLevel: number;
  quizCount: number;
  correctCount: number
}

export interface QuizSettingsType {
  quizSettingCompleted: boolean;
  questionCount: number;
  isSaveResult: boolean;
  isRegisteredWords: number;
  correctRate: number;
  attentionLevelList: number[];
  partsOfSpeeches: number[];
  isIdioms: number;
  isSpecialCharacters: number;
}