import React from 'react'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import MyPage from './MyPage'
import { User } from '../../types/userTypes'

// テスト用のユーザー情報を用意
const mockUser: User = {
  name: 'Test User',
  admin: false,
  root: false,
}

describe('MyPage Component', () => {
  test('ユーザー名が表示される', () => {
    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    // ユーザー名
    expect(screen.getByText('ようこそ、Test Userさん！')).toBeInTheDocument()
  })

  test('管理ユーザーの場合にメッセージが表示される', () => {
    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    // 管理ユーザー用メッセージ
    expect(
      screen.getByText('管理ユーザーでログインしています。'),
    ).toBeInTheDocument()
  })

  test('通常ユーザーの場合は管理ユーザー用メッセージが表示されない', () => {
    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    // 管理ユーザー用メッセージが表示されないことを確認
    expect(
      screen.queryByText('管理ユーザーでログインしています。'),
    ).not.toBeInTheDocument()
  })

  test('今日の日付が表示される', () => {
    // テスト日時を固定にしたい場合は jest.useFakeTimers などでDateをモック化できる
    // ここでは単純に「部分一致するか」だけをテスト
    const todayString = new Date().toLocaleDateString()

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    expect(
      screen.getByText((content) => content.includes(todayString)),
    ).toBeInTheDocument()
  })

  test('サインアウトボタンがクリックされたら onSignOut が呼ばれる', async () => {
    const onSignOutMock = jest.fn()

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    // サインアウトボタンをクリック
    await userEvent.click(screen.getByRole('button', { name: 'サインアウト' }))

    expect(onSignOutMock).toHaveBeenCalledTimes(1)
  })
})
