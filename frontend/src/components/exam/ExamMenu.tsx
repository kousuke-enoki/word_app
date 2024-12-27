import React, { useState } from 'react'
import ExamOptions from './ExamOptions'
import ExamSettings from './ExamSettings'
import ExamStart from './ExamStart'

type ExamType = 'all' | 'registered' | 'custom' | 'start'

const ExamMenu: React.FC = () => {
  const [selectedExamType, setSelectedExamType] = useState<ExamType | null>(
    null,
  )
  const [customSettings, setCustomSettings] = useState({
    questionCount: 10,
    targetWordTypes: 'all',
    partsOfSpeeches: [1, 3, 4, 5] as number[],
  })

  const handleExamTypeSelect = (testType: ExamType) => {
    setSelectedExamType(testType)
    if (testType !== 'custom') {
      setCustomSettings({
        ...customSettings,
        targetWordTypes: testType,
      })
    }
  }

  const handleSaveSettings = (settings: typeof customSettings) => {
    setCustomSettings(settings)
    setSelectedExamType('custom')
  }

  return (
    <div>
      {!selectedExamType && <ExamOptions onSelect={handleExamTypeSelect} />}

      {selectedExamType === 'custom' && (
        <ExamSettings
          settings={customSettings}
          onSaveSettings={handleSaveSettings}
          onStart={() => setSelectedExamType('start')}
        />
      )}

      {selectedExamType && selectedExamType !== 'custom' && (
        <ExamStart settings={customSettings} />
      )}
    </div>
  )
}

export default ExamMenu
