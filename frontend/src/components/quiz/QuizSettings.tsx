import React, { useState } from 'react';
import { getPartOfSpeech } from '../../service/word/GetPartOfSpeech';
import { QuizSettingsType } from '../../types/quiz';
import { MySwitch } from '../myUi/MySwitch';
import { MyNumberInput } from '../myUi/MyNumberInput';
import { MySelect } from '../myUi/MySelect';
import { MyCheckbox } from '../myUi/MyCheckBox';
import { MySegment } from '../myUi/MySegment';
import { MyCollapsible } from '../myUi/MyCollapsible';
import clsx from 'clsx';
import '../../styles/ui.css'

const targetOptions = [
  { value: 0, label: '全単語' },
  { value: 1, label: '登録単語のみ' },
];
const attentionLevels = [1,2,3,4,5]; 

type QuizSettingsProps = {
  settings: QuizSettingsType
  onSaveSettings: (v: QuizSettingsType) => void;
};

const QuizSettings: React.FC<QuizSettingsProps> = ({
  settings,
  onSaveSettings,
}) => {
  const [localSettings, setLocalSettings] = useState<QuizSettingsType>(settings);

  const handleSave = () => {
    console.log(localSettings)

    onSaveSettings({ ...localSettings, quizSettingCompleted: true });
  };
  const isRegisteredMode = localSettings.isRegisteredWords === 1;

  const upd = (k: keyof QuizSettingsType, v: any)=> setLocalSettings(p=>({...p,[k]:v}));

  return (
    <div className="ui-card max-w-lg mx-auto">
      <h2 className="text-lg font-bold mb-4">テスト設定</h2>

      {/* ── 必須項目 ────────────────────────────── */}
      <section className="ui-section">
        <span className="ui-heading">必須</span>

        {/* 問題数 */}
        <div className="flex items-center gap-3">
          <label className="ui-label shrink-0">問題数</label>
          <MyNumberInput value={localSettings.questionCount} min={10} max={100}
                         onChange={v=>upd('questionCount',v)} />
          <span className="text-sm">問</span>
        </div>

        {/* 対象単語 */}
        <div className="flex items-center gap-3">
          <label className="ui-label shrink-0">対象</label>
          <MySegment value={localSettings.isRegisteredWords}
                     onChange={v=>upd('isRegisteredWords',Number(v))}
                     targets={targetOptions}/>
        </div>

        {/* 成績保存 */}
        <div className="flex items-center gap-3">
          <label className="ui-label shrink-0">成績保存</label>
          <MySwitch id="saveRes" checked={localSettings.isSaveResult}
                    onChange={v=>upd('isSaveResult',v)} />
        </div>
      </section>

      {/* ── 登録単語用オプション ───────────────── */}
      <section className={clsx("ui-section", !isRegisteredMode && "ui-disabled")}>
        <span className="ui-heading">登録単語オプション</span>

        {/* 正解率 */}
          <MyCollapsible title="正解率" disabled={!isRegisteredMode}>
          <div className="flex items-center gap-2">
            <MyNumberInput value={localSettings.correctRate} min={0} max={100}
                           onChange={v=>upd('correctRate',v)} />
            <span className="text-sm">% 以下</span>
          </div>
        </MyCollapsible>

        {/* 注意レベル */}
          <MyCollapsible title="注意レベル" disabled={!isRegisteredMode}>
          <div className="grid grid-cols-5 gap-2">
            {attentionLevels.map(l=>(
              <MyCheckbox key={l}
                label={l+''}
                checked={localSettings.attentionLevelList.includes(l)}
                onChange={()=>upd('attentionLevelList',
                  localSettings.attentionLevelList.includes(l)
                    ? localSettings.attentionLevelList.filter(x=>x!==l)
                    : [...localSettings.attentionLevelList,l])}/>
            ))}
          </div>
        </MyCollapsible>
      </section>

      <section className="ui-section">
        <span className="ui-heading">その他</span>

        {/* 品詞 */}
        <MyCollapsible title="出題する品詞" disabled={false}>
          <div className="grid grid-cols-2 gap-2">
            {getPartOfSpeech.map(p=>(
              <MyCheckbox key={p.id}
                label={p.name}
                checked={localSettings.partsOfSpeeches.includes(p.id)}
                onChange={()=>upd('partsOfSpeeches',
                  localSettings.partsOfSpeeches.includes(p.id)
                    ? localSettings.partsOfSpeeches.filter(x=>x!==p.id)
                    : [...localSettings.partsOfSpeeches,p.id])}/>
            ))}
          </div>
        </MyCollapsible>

        {/* 慣用句 / 特殊文字 */}
        <div className="flex items-center gap-3">
          <label className="ui-label">慣用句</label>
          <MySelect options={targetOptions} value={localSettings.isIdioms}
                    onChange={v=>upd('isIdioms',v)}/>
        </div>
        <div className="flex items-center gap-3">
          <label className="ui-label">特殊文字</label>
          <MySelect options={targetOptions} value={localSettings.isSpecialCharacters}
                    onChange={v=>upd('isSpecialCharacters',v)}/>
        </div>
      </section>

      {/* ── ボタン ───────────────────────────── */}
      <button 
        onClick={handleSave}
              className="mt-6 w-full rounded-md bg-primary py-2 text-white font-semibold
                         hover:bg-primary/90 transition">
        上記の設定でテスト開始
      </button>
    </div>
  );
};

export default QuizSettings;
