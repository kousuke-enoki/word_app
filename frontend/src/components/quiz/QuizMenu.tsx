import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import axiosInstance from '../../axiosConfig'
import { AnswerRouteRes, ChoiceJpm, QuizQuestion, QuizSettingsType } from '../../types/quiz';
import QuizQuestionView from './QuizQuestionView'
import QuizSettings from './QuizSettings'
import QuizStart from './QuizStart'

type QuizState = 'pause' | 'setting' | 'create' | 'running'

const QuizMenu: React.FC = () => {
  const [quizState, setQuizState] = useState<QuizState | null>(
    'pause',
  )
  const [quizSettings, setQuizSettings] = useState<QuizSettingsType>({
    quizSettingCompleted: false,
    questionCount: 10,
    isSaveResult: true,
    isRegisteredWords: 0,
    correctRate: 100,
    attentionLevelList: [1, 2, 3, 4, 5] as number[],
    partsOfSpeeches: [1, 3, 4, 5] as number[],
    isIdioms: 0,
    isSpecialCharacters: 0,
  })
  const [loading, setLoading] = useState<boolean>(true)
  const [question, setQuestion] = useState<QuizQuestion|null>(null)
  const nav = useNavigate()
  useEffect(() => {
    const fetchQuiz = async () => {
      try {
        const response = await axiosInstance.get(`/quizzes`)
        if(response.data.isRunningQuiz){
          const { nextQuestion } = response.data;
          const quizQuestion: QuizQuestion = {
            quizID: nextQuestion.quizID,
            questionNumber: nextQuestion.questionNumber,
            wordName: nextQuestion.wordName,
            choicesJpms: nextQuestion.choicesJpms.map((c:ChoiceJpm) => ({
              japaneseMeanID: c.japaneseMeanID,
              name: c.name,
            })),
          };
          setQuestion(quizQuestion)
          setQuizState('running')
        } else {
          setQuizState('setting')
        }
      } catch (error) {
        console.error(error)
        setQuizState('setting')
      } finally {
        setLoading(false)
      }
    }
    fetchQuiz()
  }, [])

  const handleSaveSettings = (settings: QuizSettingsType) => {
    setQuizSettings(settings)
    setQuizState('create')
  }

  const handleQuizCreated = (id: number, first: QuizQuestion) => {
    setQuestion(first);
    setQuizState('running');
  };
  
  const handleQuizCreateError = (msg: string) => {
    alert(msg);
    setQuizState('setting');
  };

  const handleAnswerRoute = (res: AnswerRouteRes) => {
    if (res.isFinish == false && res.nextQuestion) {
      setQuestion(res.nextQuestion);
    } else if (res.isFinish == true) {
      nav(`/results/${res.quizNumber}`)
    }
  };
  
  return (
    <div>
      {loading || quizState == 'pause' &&(
        <p>Loading...</p>
      )}

      {!quizState || quizState == 'setting' && (
        <QuizSettings
          settings={quizSettings}
          onSaveSettings={handleSaveSettings}
        />
      )}

      {quizState === 'create' && (
        <QuizStart
          settings={quizSettings}
          onSuccess={handleQuizCreated}
          onFail={handleQuizCreateError}
        />
      )}

      {quizState === 'running' && question && (
        <QuizQuestionView question={question} onAnswered={handleAnswerRoute} />
      )}
    </div>
  )
}

export default QuizMenu
