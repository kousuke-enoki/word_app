import React, { useState } from 'react'
import QuizOptions from './QuizOptions'
import QuizSettings from './QuizSettings'
import QuizStart from './QuizStart'

type QuizType = 'all' | 'registered' | 'custom' | 'start'

const QuizMenu: React.FC = () => {
  const [selectedQuizType, setSelectedQuizType] = useState<QuizType | null>(
    null,
  )
  const [customSettings, setCustomSettings] = useState({
    questionCount: 10,
    isSaveResult: true,
    targetWordTypes: 'all',
    partsOfSpeeches: [1, 3, 4, 5] as number[],
  })
  console.log("asdf")
  const handleQuizTypeSelect = (testType: QuizType) => {
    setSelectedQuizType(testType)
    if (testType !== 'custom') {
      setCustomSettings({
        ...customSettings,
        targetWordTypes: testType,
      })
    }
  }

  const handleSaveSettings = (settings: typeof customSettings) => {
    setCustomSettings(settings)
    setSelectedQuizType('custom')
  }

  return (
    <div>
    {'quiz menu'}
      {!selectedQuizType && <QuizOptions onSelect={handleQuizTypeSelect} />}

      {selectedQuizType === 'custom' && (
        <QuizSettings
          settings={customSettings}
          onSaveSettings={handleSaveSettings}
          onStart={() => setSelectedQuizType('start')}
        />
      )}

      {selectedQuizType && selectedQuizType !== 'custom' && (
        <QuizStart settings={customSettings} />
      )}
    </div>
  )
}

export default QuizMenu
