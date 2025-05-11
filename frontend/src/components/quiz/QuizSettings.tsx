import React, { useState } from 'react';
import { getPartOfSpeech } from '../../service/word/GetPartOfSpeech';
import { QuizSettingsType } from '../../types/quiz';
import { MySwitch } from '../myUi/MySwitch';
import { MyNumberInput } from '../myUi/MyNumberInput';
import { MySelect } from '../myUi/MySelect';
import { MyCheckbox } from '../myUi/MyCheckBox';
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

const targetOptions = [
  { value: 0, label: '全単語' },
  { value: 1, label: '登録単語のみ' },
  // { value: 2, label: '未登録単語のみ' },
];

type QuizSettingsProps = any;

const QuizSettings: React.FC<QuizSettingsProps> = ({
  settings,
  onSaveSettings,
}) => {
  const [localSettings, setLocalSettings] = useState<QuizSettingsType>(settings);
  // const partsOfSpeechForQuiz = getPartOfSpeech;

  const handleSave = () => {
    onSaveSettings({ ...localSettings, quizSettingCompleted: true });
  };

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-bold">テスト設定</h2>

      {/* 問題数 */}
      <div className="flex items-center gap-2">
        <label htmlFor="q-num">問題数:</label>
        <MyNumberInput
          value={localSettings.questionCount}
          min={10}
          max={100}
          onChange={(v) =>
            setLocalSettings((p:any) => ({ ...p, questionCount: v }))
          }
        />{'問 出題'}
      </div>

      {/* 成績保存スイッチ */}
      <div className="flex items-center gap-3">
        <label htmlFor="save-result-switch">成績を保存:</label>
        <MySwitch
          id="save-result-switch"
          checked={localSettings.isSaveResult}
          onChange={(v) =>
            setLocalSettings((prev:any) => ({ ...prev, isSaveResult: v }))
          }
        />
      </div>

      {/* 対象単語 */}
      <div className="flex items-center gap-3">
        <label>登録単語:</label>
        <MySelect
          options={targetOptions}
          value={localSettings.isRegisteredWords}
          onChange={(v) =>
            setLocalSettings((p:any) => ({ ...p, isRegisteredWords: v }))
          }
        />
      </div>

      {/* 正解率 */}
      <div className="flex items-center gap-2">
        <label htmlFor="q-num">正解率:</label>
        <MyNumberInput
          value={localSettings.correctRate}
          min={0}
          max={100}
          onChange={(v) =>
            setLocalSettings((p:any) => ({ ...p, correctRate: v }))
          }
        />{'% 以下のものを出題'}
      </div>

      {/* 注意レベル */}
      <div>
        <label className="block mb-1">注意レベル:</label>
        <div className="grid grid-cols-2 gap-2">
          {localSettings.attentionLevelList.map((attentionLevel: number) => (
            <MyCheckbox
              key={attentionLevel}
              checked={localSettings.attentionLevelList.includes(attentionLevel)}
              label={attentionLevel.toString()}
              onChange={() => {
                setLocalSettings((prev: any) => {
                  const list = prev.attentionLevelList.includes(attentionLevel)
                    ? prev.attentionLevelList.filter((id: number) => id !== attentionLevel)
                    : [...prev.attentionLevelList, attentionLevel];
                  return { ...prev, attentionLevelList: list };
                });
              }}
            />
          ))}
        </div>
      </div>

      {/* 品詞チェックボックス */}
      <div>
        <label className="block mb-1">出題する品詞:</label>
        <div className="grid grid-cols-2 gap-2">
          {getPartOfSpeech.map((pos) => (
            <MyCheckbox
              key={pos.id}
              checked={localSettings.partsOfSpeeches.includes(pos.id)}
              label={pos.name}
              onChange={() => {
                setLocalSettings((prev:any) => {
                  const list = prev.partsOfSpeeches.includes(pos.id)
                    ? prev.partsOfSpeeches.filter((id:number) => id !== pos.id)
                    : [...prev.partsOfSpeeches, pos.id];
                  return { ...prev, partsOfSpeeches: list };
                });
              }}
            />
          ))}
        </div>
      </div>

      {/* 対象単語 */}
      <div className="flex items-center gap-3">
        <label>慣用句を出題するかどうか:</label>
        <MySelect
          options={targetOptions}
          value={localSettings.isIdioms}
          onChange={(v) =>
            setLocalSettings((p:any) => ({ ...p, isIdioms: v }))
          }
        />
      </div>

      {/* 対象単語 */}
      <div className="flex items-center gap-3">
        <label>登録単語:</label>
        <MySelect
          options={targetOptions}
          value={localSettings.isSpecialCharacters}
          onChange={(v) =>
            setLocalSettings((p:any) => ({ ...p, isSpecialCharacters: v }))
          }
        />
      </div>

      <button
        onClick={handleSave}
        className="rounded bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700"
      >
        上記の設定でテスト開始
      </button>
    </div>
  );
};

export default QuizSettings;
