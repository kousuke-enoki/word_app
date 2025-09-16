// src/components/quiz/__tests__/QuizSettings.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import type { QuizSettingsType } from '@/types/quiz'

import QuizSettings from '../QuizSettings'

/* ================= モック ================= */

// Page 系は薄く
vi.mock('@/components/common/PageTitle', () => ({
  default: ({ title }: any) => <h1 data-testid="PageTitle">{title}</h1>,
}))
vi.mock('@/components/common/PageBottomNav', () => ({
  default: (p: any) => <div data-testid="PageBottomNav" {...p} />,
}))

// Card/Badge/Button は最小限
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
  Badge: ({ children, ...rest }: any) => (
    <span data-testid="Badge" {...rest}>
      {children}
    </span>
  ),
}))
vi.mock('@/components/ui/ui', () => ({
  Button: ({ children, onClick, className }: any) => (
    <button data-testid="Button" onClick={onClick} className={className}>
      {children}
    </button>
  ),
}))

// myUi コンポーネント郡を操作しやすいようにモック
vi.mock('@/components/myUi/MyNumberInput', () => ({
  MyNumberInput: ({ value, min, max, onChange }: any) => (
    <input
      data-testid="MyNumberInput"
      type="number"
      value={value}
      min={min}
      max={max}
      onChange={(e) => onChange(Number(e.target.value))}
    />
  ),
}))
vi.mock('@/components/myUi/MySegment', () => ({
  MySegment: ({ value, targets, onChange }: any) => (
    <div data-testid="MySegment" data-current={value}>
      {targets.map((t: any) => (
        <button
          key={t.value}
          aria-pressed={value === t.value}
          onClick={() => onChange(t.value)}
        >
          {t.label}
        </button>
      ))}
    </div>
  ),
}))
vi.mock('@/components/myUi/MySwitch', () => ({
  MySwitch: ({ checked, onChange, id }: any) => (
    <label>
      <input
        data-testid={id || 'MySwitch'}
        type="checkbox"
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
      />
    </label>
  ),
}))
vi.mock('@/components/myUi/MyCollapsible', () => ({
  MyCollapsible: ({
    title,
    disabled,
    defaultOpen,
    children,
    removeCollapsedGap,
  }: any) => (
    <section
      data-testid={`MyCollapsible:${title}`}
      data-disabled={String(!!disabled)}
      data-open={String(!!defaultOpen)}
      data-remove-gap={String(!!removeCollapsedGap)}
    >
      <h3>{title}</h3>
      {/* モックでは常に children を表示 */}
      <div>{children}</div>
    </section>
  ),
}))
vi.mock('@/components/myUi/MyCheckBox', () => ({
  MyCheckbox: ({ label, checked, onChange }: any) => (
    <label>
      <input
        data-testid={`MyCheckbox:${label}`}
        type="checkbox"
        checked={checked}
        onChange={onChange}
      />
      <span>{label}</span>
    </label>
  ),
}))
vi.mock('@/components/myUi/MySelect', () => ({
  MySelect: ({ options, value, onChange }: any) => (
    <select
      data-testid="MySelect"
      value={value}
      onChange={(e) => onChange(Number(e.target.value))}
    >
      {options.map((o: any) => (
        <option key={o.value} value={o.value}>
          {o.label}
        </option>
      ))}
    </select>
  ),
}))

// 品詞一覧（配列）をモック。id=1,3,4,5 を含める
vi.mock('@/service/word/GetPartOfSpeech', () => ({
  getPartOfSpeech: [
    { id: 1, name: '名詞' },
    { id: 2, name: '代名' },
    { id: 3, name: '動詞' },
    { id: 4, name: '形容' },
    { id: 5, name: '副詞' },
  ],
}))

/* ================= ヘルパ ================= */

const base: QuizSettingsType = {
  quizSettingCompleted: false,
  questionCount: 10,
  isSaveResult: true,
  isRegisteredWords: 0,
  correctRate: 100,
  attentionLevelList: [1, 2, 3, 4, 5],
  partsOfSpeeches: [1, 3, 4, 5],
  isIdioms: 0,
  isSpecialCharacters: 0,
}

const renderWith = (settings: QuizSettingsType, onSave = vi.fn()) =>
  render(<QuizSettings settings={settings} onSaveSettings={onSave} />)

beforeEach(() => {
  vi.clearAllMocks()
})

/* ================= テスト本体 ================= */

describe('QuizSettings', () => {
  it('初期表示：タイトル/バッジ/必須項目/その他のラベルとデフォルト値', () => {
    renderWith(base)

    // 見出しとバッジ
    expect(screen.getByTestId('PageTitle')).toHaveTextContent('テスト設定')
    expect(screen.getByTestId('Badge')).toHaveTextContent('クイズ')

    // 問題数
    const nums = screen.getAllByTestId('MyNumberInput')
    expect(nums[0]).toHaveValue(10) // 最初の NumberInput は questionCount

    // 対象（セグメント）
    const seg = screen.getByTestId('MySegment')
    expect(seg).toHaveAttribute('data-current', '0') // 全単語

    // 成績保存
    expect(screen.getByTestId('saveRes')).toBeChecked()

    // 「登録単語オプション」セクションは isRegisteredWords=0 で半透明（クラス付与）
    // const regOptSection =
    //   screen
    //     .getByText('登録単語オプション')
    //     .closest('section')!
    //     .parentElement!.closest('section') || // h3→MyCollapsible→section 構造のため安全側に辿る
    //   screen.getByText('登録単語オプション').closest('section')
    // 親 Card 内に opacity-60 がある（実装では <section className={clsx(!isRegisteredMode && 'opacity-60')}>）
    // モックでは className がそのまま data-testid="Card" 上に乗るので、title 付近から探す
    const card = screen.getAllByTestId('Card')[0]
    expect(card.innerHTML).toContain('登録単語オプション')

    // 品詞チェックボックス（モック=5個）
    const noun = screen.getByTestId('MyCheckbox:名詞')
    const pron = screen.getByTestId('MyCheckbox:代名')
    const verb = screen.getByTestId('MyCheckbox:動詞')
    const adj = screen.getByTestId('MyCheckbox:形容')
    const adv = screen.getByTestId('MyCheckbox:副詞')

    expect(noun).toBeChecked() // 初期 partsOfSpeeches に 1
    expect(pron).not.toBeChecked() // 2 は含まれない
    expect(verb).toBeChecked() // 3
    expect(adj).toBeChecked() // 4
    expect(adv).toBeChecked() // 5

    // セレクトのデフォルト値
    const selects = screen.getAllByTestId('MySelect')
    // 0: 慣用句 / 1: 特殊文字
    expect(selects[0]).toHaveValue('0')
    expect(selects[1]).toHaveValue('0')
  })

  it('入力変更 → 保存で onSaveSettings に quizSettingCompleted=true 付きで渡る', async () => {
    const onSave = vi.fn()
    renderWith(base, onSave)

    // 問題数: 20
    const numbers = screen.getAllByTestId('MyNumberInput')
    await userEvent.clear(numbers[0])
    await userEvent.type(numbers[0], '20')

    // 対象: 登録単語のみ（value=1）
    const seg = screen.getByTestId('MySegment')
    await userEvent.click(
      within(seg).getByRole('button', { name: '登録単語のみ' }),
    )

    // 成績保存: OFF
    await userEvent.click(screen.getByTestId('saveRes'))

    // 正解率: 80
    const rateInput = screen.getAllByTestId('MyNumberInput')[1]
    await userEvent.clear(rateInput)
    await userEvent.type(rateInput, '80')

    // 注意レベル: 「2」と「5」を外す
    await userEvent.click(screen.getByTestId('MyCheckbox:2'))
    await userEvent.click(screen.getByTestId('MyCheckbox:5'))

    // 品詞: 名詞(1)を外し、代名(2)を追加
    await userEvent.click(screen.getByTestId('MyCheckbox:名詞'))
    await userEvent.click(screen.getByTestId('MyCheckbox:代名'))

    // 慣用句: 含まない(2)
    const selects = screen.getAllByTestId('MySelect')
    await userEvent.selectOptions(selects[0], '2')
    // 特殊文字: 含む(1)
    await userEvent.selectOptions(selects[1], '1')

    // 保存
    await userEvent.click(
      screen.getByRole('button', { name: '上記の設定でテスト開始' }),
    )

    expect(onSave).toHaveBeenCalledTimes(1)
    const payload = onSave.mock.calls[0][0] as QuizSettingsType

    // 変化を検証
    expect(payload.quizSettingCompleted).toBe(true)
    expect(payload.questionCount).toBe(20)
    expect(payload.isRegisteredWords).toBe(1)
    expect(payload.isSaveResult).toBe(false)
    expect(payload.correctRate).toBe(80)
    expect(new Set(payload.attentionLevelList)).toEqual(new Set([1, 3, 4])) // 2,5 を外した
    expect(new Set(payload.partsOfSpeeches)).toEqual(new Set([2, 3, 4, 5])) // 1→外し、2→追加
    expect(payload.isIdioms).toBe(2)
    expect(payload.isSpecialCharacters).toBe(1)
  })

  it('登録単語モードの切り替えで「登録単語オプション」Collapsible の disabled/defaultOpen が反映される', async () => {
    renderWith(base) // isRegisteredWords=0

    // まずは disabled=true, defaultOpen=false（モックは defaultOpen を data-open で表示）
    const rateColl = screen.getByTestId('MyCollapsible:正解率')
    const attnColl = screen.getByTestId('MyCollapsible:注意レベル')
    expect(rateColl).toHaveAttribute('data-disabled', 'true')
    expect(attnColl).toHaveAttribute('data-disabled', 'true')

    // モード ON（登録単語のみ=1）
    const seg = screen.getByTestId('MySegment')
    await userEvent.click(
      within(seg).getByRole('button', { name: '登録単語のみ' }),
    )

    // 再取得（DOM 再描画後の属性）
    const rateColl2 = screen.getByTestId('MyCollapsible:正解率')
    const attnColl2 = screen.getByTestId('MyCollapsible:注意レベル')
    expect(rateColl2).toHaveAttribute('data-disabled', 'false')
    expect(attnColl2).toHaveAttribute('data-disabled', 'false')
    // 正解率は defaultOpen が true の仕様
    expect(rateColl2).toHaveAttribute('data-open', 'true')
  })

  it.each([
    { idioms: 0, special: 0 },
    { idioms: 0, special: 1 },
    { idioms: 0, special: 2 },
    { idioms: 1, special: 0 },
    { idioms: 1, special: 1 },
    { idioms: 1, special: 2 },
    { idioms: 2, special: 0 },
    { idioms: 2, special: 1 },
    { idioms: 2, special: 2 },
  ])(
    '慣用句/特殊文字の全パターン (idioms=$idioms, special=$special)',
    async ({ idioms, special }) => {
      const onSave = vi.fn()
      renderWith(base, onSave)

      const selects = screen.getAllByTestId('MySelect')
      await userEvent.selectOptions(selects[0], String(idioms))
      await userEvent.selectOptions(selects[1], String(special))

      await userEvent.click(
        screen.getByRole('button', { name: '上記の設定でテスト開始' }),
      )
      const payload = onSave.mock.calls[0][0] as QuizSettingsType
      expect(payload.isIdioms).toBe(idioms)
      expect(payload.isSpecialCharacters).toBe(special)
    },
  )

  it('注意レベルと品詞のトグルは複数回押下で add/remove を繰り返す', async () => {
    const onSave = vi.fn()
    renderWith(base, onSave)

    // 注意レベル 3 を外す→戻す
    const att3 = screen.getByTestId('MyCheckbox:3')
    await userEvent.click(att3) // 外す
    await userEvent.click(att3) // 戻す

    // 品詞: 動詞(3)を外す→戻す
    const verb = screen.getByTestId('MyCheckbox:動詞')
    await userEvent.click(verb) // 外す
    await userEvent.click(verb) // 戻す

    await userEvent.click(
      screen.getByRole('button', { name: '上記の設定でテスト開始' }),
    )
    const payload = onSave.mock.calls[0][0] as QuizSettingsType

    expect(new Set(payload.attentionLevelList)).toEqual(
      new Set([1, 2, 3, 4, 5]),
    )
    expect(new Set(payload.partsOfSpeeches)).toEqual(new Set([1, 3, 4, 5]))
  })
})
