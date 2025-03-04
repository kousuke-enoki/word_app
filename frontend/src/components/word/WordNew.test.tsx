// WordNew.test.tsx
import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom'
import WordNew from './WordNew' // 実際のファイルパスに変更
import axiosInstance from '../../axiosConfig'
import { useNavigate } from 'react-router-dom'

jest.mock('../../axiosConfig', () => ({
  __esModule: true,
  default: {
    post: jest.fn(),
  },
}))

// React Router の navigate をモック
jest.mock('react-router-dom', () => {
  const actual = jest.requireActual('react-router-dom')
  return {
    __esModule: true,
    ...actual,
    useNavigate: jest.fn(),
  }
})

// window.alert をモック
const mockAlert = jest.fn()
window.alert = mockAlert

describe('WordNew Component', () => {
  const mockNavigate = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    ;(useNavigate as jest.Mock).mockReturnValue(mockNavigate)
  })

  it('初期レンダリングで、単語名入力欄・品詞(1枠)・日本語訳(1枠)・登録ボタンなどが表示される', () => {
    render(<WordNew />)

    expect(screen.getByText('単語登録フォーム')).toBeInTheDocument()

    // 単語名
    const wordNameInput = screen.getByLabelText('単語名:') as HTMLInputElement
    expect(wordNameInput).toBeInTheDocument()
    expect(wordNameInput.value).toBe('')

    // 品詞 (select) が1つだけ
    const selects = screen.getAllByLabelText('品詞:')
    expect(selects.length).toBe(1)
    expect((selects[0] as HTMLSelectElement).value).toBe('0') // "選択してください"

    // 日本語訳 (input) が1つだけ
    const jpMeanInputs = screen.getAllByLabelText('日本語訳:')
    expect(jpMeanInputs.length).toBe(1)
    expect((jpMeanInputs[0] as HTMLInputElement).value).toBe('')

    // 登録ボタン
    expect(
      screen.getByRole('button', { name: '単語を登録' }),
    ).toBeInTheDocument()
    // 戻るボタン
    expect(
      screen.getByRole('button', { name: 'mypageに戻る' }),
    ).toBeInTheDocument()
  })

  it('単語名に半角英字以外を入力すると alert が表示され、入力値が反映されない', () => {
    render(<WordNew />)
    const wordNameInput = screen.getByLabelText('単語名:') as HTMLInputElement

    // OK: 'hello'
    fireEvent.change(wordNameInput, { target: { value: 'hello' } })
    expect(wordNameInput.value).toBe('hello')
    expect(mockAlert).not.toHaveBeenCalled()

    // NG: 'こんにちは'
    fireEvent.change(wordNameInput, { target: { value: 'こんにちは' } })
    expect(mockAlert).toHaveBeenCalledWith(
      '単語名は半角アルファベットのみ入力できます。',
    )
    // 値は元に戻っている
    expect(wordNameInput.value).toBe('hello')
  })

  it('日本語訳にアルファベットを入力すると alert が表示され、入力値が反映されない', () => {
    render(<WordNew />)
    const jpMeanInput = screen.getByLabelText('日本語訳:') as HTMLInputElement

    // OK: '犬'
    fireEvent.change(jpMeanInput, { target: { value: '犬' } })
    expect(jpMeanInput.value).toBe('犬')
    expect(mockAlert).not.toHaveBeenCalled()

    // NG: 'dog'
    fireEvent.change(jpMeanInput, { target: { value: 'dog' } })
    expect(mockAlert).toHaveBeenCalledWith(
      '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。',
    )
    // 値は元の '犬' に戻る
    expect(jpMeanInput.value).toBe('犬')
  })

  it('品詞を追加・削除できる', () => {
    render(<WordNew />)
    // 初期は品詞1つ
    expect(screen.getAllByLabelText('品詞:').length).toBe(1)

    // 追加
    fireEvent.click(screen.getByRole('button', { name: '品詞を追加' }))
    expect(screen.getAllByLabelText('品詞:').length).toBe(2)

    // 2つ目に削除ボタンがついている
    const removeButtons = screen.getAllByRole('button', { name: '品詞を削除' })
    expect(removeButtons.length).toBe(2)

    // 1つ削除すると品詞は1つに戻る
    fireEvent.click(removeButtons[1])
    expect(screen.getAllByLabelText('品詞:').length).toBe(1)
  })

  it('日本語訳を追加・削除できる', () => {
    render(<WordNew />)
    // 初期は日本語訳1つ
    expect(screen.getAllByLabelText('日本語訳:').length).toBe(1)

    // 追加
    fireEvent.click(screen.getByRole('button', { name: '日本語訳を追加' }))
    expect(screen.getAllByLabelText('日本語訳:').length).toBe(2)

    // 2つ目の削除ボタンを押す
    const removeMeanButtons = screen.getAllByRole('button', { name: '削除' })
    fireEvent.click(removeMeanButtons[1])
    expect(screen.getAllByLabelText('日本語訳:').length).toBe(1)
  })

  it('バリデーションチェックに引っかかる場合は送信せずエラーメッセージが表示される', async () => {
    render(<WordNew />)

    // まだ名前未入力 or 不正なまま
    // 品詞: デフォルトで0(未選択)
    // 日本語訳: '日本語' 項目としてはまだ valid とみなされる可能性(正規表現に合致しているかどうか)

    // 送信
    fireEvent.click(screen.getByRole('button', { name: '単語を登録' }))

    // validateWord でエラーが出るので axios は呼ばれない
    await waitFor(() => {
      expect(axiosInstance.post).not.toHaveBeenCalled()
    })

    // エラーメッセージ（品詞を選択してください。など）
    expect(
      screen.getByText('単語名は半角アルファベットのみ入力できます。'),
    ).toBeInTheDocument()
    expect(screen.getByText('品詞を選択してください。')).toBeInTheDocument()
  })

  it('正しく入力した場合、axiosInstance.post が呼ばれて onSuccess で navigate が呼ばれる', async () => {
    // axios のモック成功レスポンス
    ;(axiosInstance.post as jest.Mock).mockResolvedValue({
      data: { id: 123, name: 'apple' },
    })

    render(<WordNew />)

    // 単語名入力
    fireEvent.change(screen.getByLabelText('単語名:'), {
      target: { value: 'apple' },
    })
    // 品詞を0以外に
    const partOfSpeechSelect = screen.getByLabelText(
      '品詞:',
    ) as HTMLSelectElement
    // 適当に 1 or 2 など (モックの品詞がどうなっているかによる)
    fireEvent.change(partOfSpeechSelect, { target: { value: '1' } })

    // 日本語訳を「犬」に
    fireEvent.change(screen.getByLabelText('日本語訳:'), {
      target: { value: '犬' },
    })

    // 送信
    fireEvent.click(screen.getByRole('button', { name: '単語を登録' }))

    // axiosが正しく呼ばれるか
    await waitFor(() => {
      expect(axiosInstance.post).toHaveBeenCalledWith('/words/new', {
        name: 'apple',
        wordInfos: [
          {
            partOfSpeechId: 1,
            japaneseMeans: [{ name: '犬' }],
          },
        ],
      })
    })

    // onSuccess で navigate が呼ばれる
    expect(mockNavigate).toHaveBeenCalledWith('/words/123', {
      state: { successMessage: 'appleが正常に登録されました！' },
    })
  })

  it('axiosInstance.post でエラーが起きた場合、エラーメッセージが表示される', async () => {
    ;(axiosInstance.post as jest.Mock).mockRejectedValue(
      new Error('some error'),
    )

    render(<WordNew />)

    // 必要最低限の正しい入力
    fireEvent.change(screen.getByLabelText('単語名:'), {
      target: { value: 'banana' },
    })
    const partOfSpeechSelect = screen.getByLabelText(
      '品詞:',
    ) as HTMLSelectElement
    fireEvent.change(partOfSpeechSelect, { target: { value: '1' } })
    fireEvent.change(screen.getByLabelText('日本語訳:'), {
      target: { value: '猫' },
    })

    // 送信
    fireEvent.click(screen.getByRole('button', { name: '単語を登録' }))

    // エラー発生により onError
    await waitFor(() => {
      // axios呼び出し後にエラーが起きたか
      expect(axiosInstance.post).toHaveBeenCalledTimes(1)
      // 画面にエラーメッセージが表示される
      // eslint-disable-next-line testing-library/no-wait-for-multiple-assertions
      expect(
        screen.getByText('単語の登録中にエラーが発生しました。'),
      ).toBeInTheDocument()
    })
  })

  it('「mypageに戻る」ボタンをクリックすると navigate("/", {}) が呼ばれる', () => {
    render(<WordNew />)

    fireEvent.click(screen.getByRole('button', { name: 'mypageに戻る' }))
    expect(mockNavigate).toHaveBeenCalledWith('/', {})
  })
})
