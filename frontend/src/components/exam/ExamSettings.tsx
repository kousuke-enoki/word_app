import React, { useState } from 'react'
import { getPartsOfSpeechForExam } from '../../service/word/GetPartOfSpeech'

type ExamSettingsProps = {
  settings: {
    questionCount: number
    targetWordTypes: string
    partsOfSpeeches: number[]
  }
  onSaveSettings: (settings: {
    questionCount: number
    targetWordTypes: string
    partsOfSpeeches: number[]
  }) => void
  onStart: () => void
}

const ExamSettings: React.FC<ExamSettingsProps> = ({
  settings,
  onSaveSettings,
  onStart,
}) => {
  const [localSettings, setLocalSettings] = useState(settings)
  const partsOfSpeechForExam = getPartsOfSpeechForExam()

  // チェックボックスの変更処理
  const handleCheckboxChange = (id: number) => {
    const updatedPartsOfSpeech = localSettings.partsOfSpeeches.includes(id)
      ? localSettings.partsOfSpeeches.filter((partId) => partId !== id)
      : [...localSettings.partsOfSpeeches, id]

    setLocalSettings({
      ...localSettings,
      partsOfSpeeches: updatedPartsOfSpeech,
    })
  }

  const handleSave = () => {
    onSaveSettings(localSettings)
    onStart()
  }

  return (
    <div>
      <h2>カスタムテスト設定</h2>

      {/* 問題数 */}
      <div>
        <label>問題数: </label>
        <input
          type="number"
          min="1"
          max="50"
          value={localSettings.questionCount}
          onChange={(e) =>
            setLocalSettings({
              ...localSettings,
              questionCount: parseInt(e.target.value, 10),
            })
          }
        />
      </div>

      {/* 対象単語 */}
      <div>
        <label>対象単語: </label>
        <select
          value={localSettings.targetWordTypes}
          onChange={(e) =>
            setLocalSettings({
              ...localSettings,
              targetWordTypes: e.target.value,
            })
          }
        >
          <option value="all">全単語</option>
          <option value="registered">登録単語のみ</option>
        </select>
      </div>

      {/* 品詞チェックボックス */}
      <div>
        <label>品詞:</label>
        <div>
          {partsOfSpeechForExam.map((pos) => (
            <div key={pos.id}>
              <label>
                <input
                  type="checkbox"
                  checked={localSettings.partsOfSpeeches.includes(pos.id)}
                  onChange={() => handleCheckboxChange(pos.id)}
                />
                {pos.name}
              </label>
            </div>
          ))}
        </div>
      </div>

      <button onClick={handleSave}>設定を保存してテスト開始</button>
    </div>
  )
}

export default ExamSettings
