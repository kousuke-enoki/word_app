import React, { useEffect, useState } from 'react'
import axiosInstance from '../../axiosConfig'
import QuizSettings from './QuizSettings'
import QuizStart from './QuizStart'
import QuizQuestionView from './QuizQuestionView'
import QuizResult from './QuizResult'
import { QuizQuestion, ChoiceJpm, AnswerRouteRes, QuizSettingsType, ResultRes } from '../../types/quiz';

// type QuizType = 'all' | 'registered' | 'custom' | 'start'
type QuizState = 'pause' | 'setting' | 'create' | 'running' | 'result'
// interface QuizSettingsType {
//   quizSettingCompleted: boolean;
//   questionCount: number;
//   isSaveResult: boolean;
//   isRegisteredWords: number;
//   correctRate: number;
//   attentionLevelList: number[];
//   partsOfSpeeches: number[];
//   isIdioms: number;
//   isSpecialCharacters: number;
// }

interface Quiz {
  ID: number;
}

// interface QuizQuestion {
//   QuizID: number;
//   QuestionNumber: number;
//   WordName: string;
//   ChoicesJpms: {
//     JapaneseMeanID: number;
//     Name: string;
//   }[];
// }

// export interface AnswerRouteRes {
//   kind: 'next' | 'finish';
//   isCorrect: boolean;
//   nextQuestion?: QuizQuestion;     // kind === 'next'
//   result?: any;                    // kind === 'finish'
// }


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
  const [result, setResult] = useState<ResultRes|null>(null)
  console.log("asdf")
  useEffect(() => {
    const fetchQuiz = async () => {
      console.log(quizState)
      try {
        const response = await axiosInstance.get(`/quizzes`)
        console.log(response)
        if(response.data.isRunningQuiz){
          const { nextQuestion } = response.data;
          console.log(nextQuestion)
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
          console.log('if')
          console.log(question)
          setQuizState('running')
        } else {
          console.log('else')
          setQuizState('setting')
        }
        // setMemo(response.data.memo || '')
      } catch (error) {
        console.log('err')
        // alert('クイズの取得中にエラーが発生しました。')
        setQuizState('setting')
      } finally {
        console.log('finally')
        console.log(quizState)
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

  const handleQuizCreated = (id: number, first: QuizQuestion) => {
    console.log('first')
    console.log(first)
    console.log(id)
    setQuestion(first);
    setQuizState('running');
  };
  
  const handleQuizCreateError = (msg: string) => {
    console.log(msg)
    alert(msg);
    setQuizState('setting');
  };

  const handleAnswerRoute = (res: AnswerRouteRes) => {
    if (res.isFinish == false && res.nextQuestion) {
      console.log('isfinish')
      console.log(res)
      setQuestion(res.nextQuestion);
    } else if (res.isFinish == true && res.result) {
      console.log(res)
      // const { quizNumber, totalQuestionsCount, correctCount, 
      //   resultCorrectRate,resultSetting, resultQuestions  } = res.result;
      // const result: ResultRes = {
      //   quizNumber: quizNumber,
      //   totalQuestionsCount: totalQuestionsCount,
      //   correctCount: correctCount,
      //   resultCorrectRate: resultCorrectRate,
      //   resultSetting: resultSetting,
      //   resultQuestions: resultQuestions.map((rq) => ({
      //     quizID: rq.quizID,
      //     questionNumber: rq.questionNumber,
      //     wordID: rq.wordID,
      //     wordName: rq.wordName,
      //     posID: rq.posID,
      //     correctJpmId: rq.correctJpmId,
      //     choicesJpms: rq.choicesJpms.map((c) => ({
      //       japaneseMeanID: c.japaneseMeanID,
      //       name: c.name,
      //     })),
      //     answerJpmId: rq.answerJpmId,
      //     isCorrect: rq.isCorrect,
      //     timeMs: rq.timeMs,
      //     registeredWord: {
      //       isRegistered: rq.registeredWord.isRegistered,
      //       attentionLevel: rq.registeredWord.attentionLevel,
      //       quizCount: rq.registeredWord.quizCount,
      //       correctCount: rq.registeredWord.correctCount,
      //     },
      //   })),
      // };
      console.log(res.result)
      setResult(res.result);
      setQuizState('result');
    }
  };
  
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

      {quizState === 'result' && result && (
        <QuizResult result={result} />
      )}
    </div>
  )
}

export default QuizMenu
