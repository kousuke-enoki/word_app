/* eslint-disable @typescript-eslint/no-explicit-any */
const navigateMock = vi.fn();

/** 必ず「他の import より前」で行う **/
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return { ...actual, useNavigate: () => navigateMock };
});


import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'

/* -------------------- モック -------------------- */
// axiosInstance.get を好きなレスポンスに差し替えられるようにする
vi.mock('@/axiosConfig', () => ({
  default: {
    get: vi.fn(),
  },
}))

import MyPage from '../MyPage'
import axiosInstance from '@/axiosConfig'

/* -------------------- 共通セットアップ -------------------- */
beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
})

/* -------------------- テスト本体 -------------------- */
describe('MyPage Component', () => {
  it('通常ユーザーの場合、ユーザー名だけが表示される', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { user: { id: 1, name: 'Test User', isAdmin: false, isRoot: false } },
    })

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    // fetch → state 更新を待つ
    expect(
      await screen.findByText('ようこそ、Test Userさん！'),
    )

    // 管理/ルート用メッセージは出ない
    expect(screen.queryByText('管理ユーザーでログインしています。')).toBeNull()
    expect(screen.queryByText('ルートユーザーでログインしています。')).toBeNull()
  })

  it('管理ユーザーには管理メッセージとリンクが表示される', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { user: { id: 2, name: 'Admin', isAdmin: true, isRoot: false } },
    })

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    expect(
      await screen.findByText('管理ユーザーでログインしています。'),
    )
    expect(screen.getByRole('link', { name: '単語登録画面' }))
  })

  it('root ユーザーには root メッセージと root 用リンクが表示される', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { user: { id: 3, name: 'Root', isAdmin: false, isRoot: true } },
    })

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    expect(
      await screen.findByText('ルートユーザーでログインしています。'),
    )
    expect(screen.getByRole('link', { name: '管理設定画面' }))
  })
  it('サインアウトで token が消え、トップへ navigate', async () => {
    (axiosInstance.get as any).mockResolvedValueOnce({
      data: { user: { id: 1, name: 'Test', isAdmin: false, isRoot: false } },
    });
    localStorage.setItem('token', 'dummy');

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    );

    await screen.findByText('ようこそ、Testさん！');

    await userEvent.click(screen.getByRole('button', { name: 'サインアウト' }));

    await waitFor(() => {
      expect(navigateMock).toHaveBeenCalledWith('/');
    });

    expect(localStorage.getItem('token')).toBeNull();
    expect(localStorage.getItem('logoutMessage')).toBe('ログアウトしました');
  });

  it('認証エラー時に token を削除し 2 秒後にトップへリダイレクト', async () => {
    (axiosInstance.get as any).mockRejectedValueOnce(new Error('401'));
    localStorage.setItem('token', 'expired-token');
  
    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>
    );
    await screen.findByText('ユーザー情報がありません。');
    expect(localStorage.getItem('token')).toBeNull();
    expect(localStorage.getItem('logoutMessage')).toBe('ログインしてください');
  });
})
