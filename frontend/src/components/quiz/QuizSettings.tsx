/* eslint-disable @typescript-eslint/no-explicit-any */
import clsx from 'clsx'
import React, { useState } from 'react'

import { Badge, Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { getPartOfSpeech } from '@/service/word/GetPartOfSpeech'
import { QuizSettingsType } from '@/types/quiz'

import PageBottomNav from '../common/PageBottomNav'
import PageTitle from '../common/PageTitle'
import { MyCheckbox } from '../myUi/MyCheckBox'
import { MyCollapsible } from '../myUi/MyCollapsible'
import { MyNumberInput } from '../myUi/MyNumberInput'
import { MySegment } from '../myUi/MySegment'
import { MySelect } from '../myUi/MySelect'
import { MySwitch } from '../myUi/MySwitch'

const targetOptions = [
  { value: 0, label: '全単語' },
  { value: 1, label: '登録単語のみ' },
]
const idiomsTargetOptions = [
  { value: 0, label: '全て' },
  { value: 1, label: '含む' },
  { value: 2, label: '含まない' },
]
const isSpecialCharactersTargetOptions = [
  { value: 0, label: '全て' },
  { value: 1, label: '含む' },
  { value: 2, label: '含まない' },
]
const attentionLevels = [1, 2, 3, 4, 5]

type Props = {
  settings: QuizSettingsType
  onSaveSettings: (v: QuizSettingsType) => void
}

const QuizSettings: React.FC<Props> = ({ settings, onSaveSettings }) => {
  const [localSettings, setLocalSettings] = useState<QuizSettingsType>(settings)
  const isRegisteredMode = localSettings.isRegisteredWords === 1
  const upd = (k: keyof QuizSettingsType, v: any) =>
    setLocalSettings((p) => ({ ...p, [k]: v }))

  const handleSave = () => {
    onSaveSettings({ ...localSettings, quizSettingCompleted: true })
  }

  return (
    <div className="mx-auto max-w-2xl">
      <div className="mb-6 flex items-center justify-between">
        <PageTitle title="テスト設定" />
        <Badge>クイズ</Badge>
      </div>

      <Card className="p-6 space-y-6">
        {/* 必須 */}
        <section>
          <h2 className="mb-3 text-sm font-semibold opacity-80">必須</h2>

          <div className="grid gap-4 sm:grid-cols-2">
            {/* 問題数 */}
            <div className="flex items-center gap-3">
              <label className="shrink-0 text-sm opacity-80">問題数</label>
              <MyNumberInput
                value={localSettings.questionCount}
                min={10}
                max={100}
                onChange={(v) => upd('questionCount', v)}
              />
              <span className="text-sm opacity-70">問</span>
            </div>

            {/* 対象単語 */}
            <div className="flex items-center gap-3">
              <label className="shrink-0 text-sm opacity-80">対象</label>
              <MySegment
                value={localSettings.isRegisteredWords}
                onChange={(v) => upd('isRegisteredWords', Number(v))}
                targets={targetOptions}
              />
            </div>

            {/* 成績保存 */}
            <div className="flex items-center gap-3">
              <label className="shrink-0 text-sm opacity-80">成績保存</label>
              <MySwitch
                id="saveRes"
                checked={localSettings.isSaveResult}
                onChange={(v) => upd('isSaveResult', v)}
              />
            </div>
          </div>
        </section>

        {/* 登録単語オプション */}
        <section className={clsx(!isRegisteredMode && 'opacity-60')}>
          <h2 className="mb-3 text-sm font-semibold opacity-80">
            登録単語オプション
          </h2>

          <MyCollapsible
            title="正解率"
            disabled={!isRegisteredMode}
            defaultOpen={isRegisteredMode}
            removeCollapsedGap
          >
            <div className="flex items-center gap-2">
              <MyNumberInput
                value={localSettings.correctRate}
                min={0}
                max={100}
                onChange={(v) => upd('correctRate', v)}
              />
              <span className="text-sm opacity-70">% 以下</span>
            </div>
          </MyCollapsible>

          <MyCollapsible
            title="注意レベル"
            disabled={!isRegisteredMode}
            removeCollapsedGap
          >
            <div className="grid grid-cols-5 gap-2">
              {attentionLevels.map((l) => (
                <MyCheckbox
                  key={l}
                  label={String(l)}
                  checked={localSettings.attentionLevelList.includes(l)}
                  onChange={() =>
                    upd(
                      'attentionLevelList',
                      localSettings.attentionLevelList.includes(l)
                        ? localSettings.attentionLevelList.filter(
                            (x) => x !== l,
                          )
                        : [...localSettings.attentionLevelList, l],
                    )
                  }
                />
              ))}
            </div>
          </MyCollapsible>
        </section>

        {/* その他 */}
        <section>
          <h2 className="mb-3 text-sm font-semibold opacity-80">その他</h2>

          <MyCollapsible title="出題する品詞" removeCollapsedGap>
            <div className="grid grid-cols-2 gap-2 md:grid-cols-3">
              {getPartOfSpeech.map((p) => (
                <MyCheckbox
                  key={p.id}
                  label={p.name}
                  checked={localSettings.partsOfSpeeches.includes(p.id)}
                  onChange={() =>
                    upd(
                      'partsOfSpeeches',
                      localSettings.partsOfSpeeches.includes(p.id)
                        ? localSettings.partsOfSpeeches.filter(
                            (x) => x !== p.id,
                          )
                        : [...localSettings.partsOfSpeeches, p.id],
                    )
                  }
                />
              ))}
            </div>
          </MyCollapsible>

          <div className="mt-3 grid gap-3 sm:grid-cols-2">
            <div className="flex items-center gap-3">
              <label className="text-sm opacity-80">慣用句</label>
              <MySelect
                options={idiomsTargetOptions}
                value={localSettings.isIdioms}
                onChange={(v) => upd('isIdioms', v)}
              />
            </div>
            <div className="flex items-center gap-3">
              <label className="text-sm opacity-80">特殊文字</label>
              <MySelect
                options={isSpecialCharactersTargetOptions}
                value={localSettings.isSpecialCharacters}
                onChange={(v) => upd('isSpecialCharacters', v)}
              />
            </div>
          </div>
        </section>

        <Button className="w-full" onClick={handleSave}>
          上記の設定でテスト開始
        </Button>
      </Card>
      <Card className="mt1 p-2">
        <PageBottomNav className="mt-1" showHome inline compact />
      </Card>
    </div>
  )
}

export default QuizSettings
