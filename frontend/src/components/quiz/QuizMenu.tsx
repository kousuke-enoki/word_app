import React, { useEffect, useState } from 'react'
import axiosInstance from '../../axiosConfig'
import QuizSettings from './QuizSettings'
import QuizStart from './QuizStart'

// type QuizType = 'all' | 'registered' | 'custom' | 'start'
type QuizState = 'pause' | 'setting' | 'create' | 'start' | 'result'
interface QuizSettingsType {
  quizSettingCompleted: boolean;
  questionCount: number;
  isSaveResult: boolean;
  isRegisteredWord: number;
  correctRate: number;
  attentionLevelList: number[];
  partsOfSpeeches: number[];
  isIdioms: number;
  isSpecialCharacters: number;
}

interface Quiz {
  ID: number;
}

interface QuizQuestion {
  QuizID: number;
  QuestionNumber: number;
  WordName: string;
  PosID: number;
  // CorrectJpmId: number;
  ChoicesJpms: ChoiceJpm[];
  // AnswerJpmId: number;
  // IsCorrect : boolean;
  // TimeMs: number;
}

interface ChoiceJpm {
  JapaneseMeanID: number;
  Name: string;
}


const QuizMenu: React.FC = () => {
  const [quizState, setQuizState] = useState<QuizState | null>(
    'pause',
  )
  const [quizSettings, setQuizSettings] = useState<QuizSettingsType>({
    quizSettingCompleted: false,
    questionCount: 10,
    isSaveResult: true,
    isRegisteredWord: 0,
    correctRate: 100,
    attentionLevelList: [1, 2, 3, 4, 5] as number[],
    partsOfSpeeches: [1, 3, 4, 5] as number[],
    isIdioms: 0,
    isSpecialCharacters: 0,
  })
  const [loading, setLoading] = useState<boolean>(true)
  const [question, setQuestion] = useState<QuizQuestion|null>(null)
  console.log("asdf")
  useEffect(() => {
    const fetchQuiz = async () => {
      try {
        const response = await axiosInstance.get(`/quizzes`)
        if(response.data.isRunningQuiz){
          setQuestion(response.data.NextQuestion)
          setQuizState('start')
        } else {
          setQuizState('setting')
        }
        // setMemo(response.data.memo || '')
      } catch (error) {
        alert('クイズの取得中にエラーが発生しました。')
      } finally {
        setLoading(false)
      }
    }
    fetchQuiz()
  }, [])
  // const handleQuizTypeSelect = (testType: QuizType) => {
  //   setSelectedQuizType(testType)
  //   if (testType !== 'custom') {
  //     setQuizSettings({
  //       ...quizSettings,
  //       targetWordTypes: testType,
  //     })
  //   }
  // }

  const handleSaveSettings = (settings: QuizSettingsType) => {
    setQuizSettings(settings)
    setQuizState('create')
  }

  return (
    <div>
    {'quiz menu'}
      {/* {!selectedQuizType && <QuizOptions onSelect={handleQuizTypeSelect} />} */}

      {loading || quizState == 'pause' &&(
        <p>Loading...</p>
      )}

      {!quizState || quizState == 'setting' && (
        <QuizSettings
          settings={quizSettings}
          onSaveSettings={handleSaveSettings}
        />
      )}

      {quizState && quizState == 'create' && (
        <QuizStart settings={quizSettings} />
      )}

      {quizState && quizState == 'start' && (
        <QuizQuestion settings={quizSettings} />
      )}
    </div>
  )
}

export default QuizMenu
