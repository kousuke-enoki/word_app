/* eslint-disable @typescript-eslint/no-explicit-any */
/* -------------------------------------------------
 * 事前：ルータ／axios／ThemeContext を全部モック
 * ------------------------------------------------- */
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import SignIn from '../SignIn'

/* ---------- react-router の useNavigate をスパイ ---------- */
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>(
    'react-router-dom',
  )
  return { ...actual, useNavigate: () => navigateMock }
})

/* ---------- axiosInstance を丸ごとモック ---------- */
vi.mock('@/axiosConfig', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))
import axiosInstance from '@/axiosConfig'

/* ---------- ThemeContext の setTheme だけ使いたい ---------- */
export const setThemeMock = vi.fn()
vi.mock('@/contexts/themeContext', () => ({
  useTheme: () => ({ setTheme: setThemeMock }),
}))

/* ---------- 各テストで state を初期化 ---------- */
beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
})
/* =================================================
 *                    TESTS
 * ================================================= */
describe('SignIn Component', () => {
  /* ---------- 1. 初期ロード ---------- */
  it('最初は Loading… → 設定取得後にフォームを表示', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: { isLineAuth: false } })

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    // スケルトン
    expect(screen.getByText('Loading…'))

    // 設定取得が終わるとフォームが現れる
    expect(
      await screen.findByRole('heading', { name: 'サインイン' }),
    )
  })

  /* ---------- 2. LINE ログインボタン ---------- */
  it('isLineAuth=true なら「LINEでログイン」ボタンが出る & クリックで location.href 書き換え', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: { is_line_auth: true } })

    // JSDOM の location を書き換え可能にする
    const originalLocation = window.location as unknown as Location;
    vi.stubGlobal('location', { ...originalLocation, href: '' });

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    const btn = await screen.findByRole('button', { name: 'LINEでログイン' })
    await userEvent.click(btn)

    expect(window.location.href).toBe(
      `${import.meta.env.VITE_API_URL}/users/auth/line/login`,
    )

    // 後片付け
    vi.stubGlobal('location', originalLocation);
  })

  /* ---------- 3. サインイン成功フロー ---------- */
  it('正しい資格情報でサインイン → token 保存・テーマ設定・/mypage へ遷移', async () => {
    // ❶ mock 順序は get → post → get の３連発
    (axiosInstance.get  as any).mockResolvedValueOnce({ data: { isLineAuth: false } });
    (axiosInstance.post as any).mockResolvedValueOnce({ data: { token: 'jwt123' } });
    (axiosInstance.get  as any).mockResolvedValueOnce({ data: { Config: { is_dark_mode: true } }, });
  
    render(<MemoryRouter><SignIn /></MemoryRouter>);
  
    /* ---------------- リアルタイマーで UI 操作 ---------------- */
    await screen.findByRole('heading', { name: 'サインイン' });
  
    const emailInput = screen.getByLabelText('Email:');
    const passInput  = screen.getByLabelText('Password:');
  
    await userEvent.type(emailInput, 'test@example.com');
    await userEvent.type(passInput,  'pass1234');
    await userEvent.click(screen.getByRole('button', { name: 'サインイン' }));
  
    /* ---------------- ここでだけフェイクタイマー ---------------- */
    vi.useFakeTimers();
    await vi.runAllTicks();     // axios の Promise を解決 → setTimeout(0) が作られる
    vi.runOnlyPendingTimers();  // setTimeout(0) を実行（navigate('/')）
    await vi.runAllTicks();     // navigate Mock が呼ばれる
    vi.useRealTimers();
  
    /* ---------------- 期待値 ---------------- */
    expect(localStorage.getItem('token')).toBe('jwt123');
    expect(setThemeMock).toHaveBeenCalledWith('dark');
    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/mypage'));
  });

  /* ---------- 4. サインイン失敗フロー ---------- */
  it('認証失敗時にはエラーメッセージを表示し token は保存しない', async () => {
    (axiosInstance.get  as any).mockResolvedValueOnce({ data: { isLineAuth: false } });
    (axiosInstance.post as any).mockRejectedValueOnce(new Error('401'));
  
    render(<MemoryRouter><SignIn /></MemoryRouter>);
  
    await screen.findByRole('heading', { name: 'サインイン' });
  
    await userEvent.type(screen.getByLabelText('Email:'),    'bad@example.com');
    await userEvent.type(screen.getByLabelText('Password:'), 'wrongpass');
    await userEvent.click(screen.getByRole('button', { name: 'サインイン' }));
  
    expect(
      await screen.findByText('Sign in failed. Please try again.')
    );
  
    expect(localStorage.getItem('token')).toBeNull();
    expect(navigateMock).not.toHaveBeenCalled();
  });
})
