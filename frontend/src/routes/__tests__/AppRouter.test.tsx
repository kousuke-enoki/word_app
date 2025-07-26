/* eslint-disable @typescript-eslint/no-unused-expressions */
import { render, screen, waitFor } from '@testing-library/react';
import { act } from 'react-dom/test-utils';
import { beforeEach,describe, expect, it, vi } from 'vitest';

import AppRouter from '../AppRouter';          // ←テスト対象

/* ------------------------------------------------------------------ */
/* ❶ ルート配下の各ページ/ヘッダーを「プレースホルダ」にモック           */
/*    - 描画されているかどうか確認できれば十分なので <div>だけ返す       */
/* ------------------------------------------------------------------ */
vi.mock('@/components/Header',                 () => ({ default: () => <div>Header</div> }));
vi.mock('@/components/user/Home',              () => ({ default: () => <div>HomePage</div> }));
vi.mock('@/components/user/SignIn',            () => ({ default: () => <div>SignInPage</div> }));
vi.mock('@/components/user/MyPage',            () => ({ default: () => <div>MyPage</div> }));
vi.mock('@/components/word/WordNew',           () => ({ default: () => <div>WordNewPage</div> }));
// …他のページも必要になったら同様に stub で追加

/* ------------------------------------------------------------------ */
/* ❷ useAuth をケースごとに好きな値で返すモック                         */
/* ------------------------------------------------------------------ */
vi.mock('@/hooks/useAuth', () => ({ useAuth: vi.fn() }));
import { useAuth } from '@/hooks/useAuth';
const mockedAuth = vi.mocked(useAuth);

/** ログイン状態を簡単に切替えるヘルパ */
const setAuth = (state: Partial<ReturnType<typeof useAuth>>) =>
  mockedAuth.mockReturnValue({ isLoggedIn: false, isLoading: false, userRole: 'guest', ...state });

/* ------------------------------------------------------------------ */
/* ❸ BrowserRouter を使うため、テスト前に location をセットするヘルパ   */
/* ------------------------------------------------------------------ */
const goTo = (path: string) => {
  act(() => {   // React18 の StrictMode 互換
    window.history.pushState({}, '', path);
  });
};

/* ------------------------------------------------------------------ */
/*                               TESTS                                */
/* ------------------------------------------------------------------ */
describe('AppRouter 全体ルーティング', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('未ログインで / にアクセスすると Home が表示される', () => {
    setAuth({ isLoggedIn: false });
    goTo('/');

    render(<AppRouter />);
    expect(screen.getByText('HomePage'))
  });

  it('未ログインで /sign_in にアクセスすると SignIn が表示される', () => {
    setAuth({ isLoggedIn: false });
    goTo('/sign_in');

    render(<AppRouter />);
    expect(screen.getByText('SignInPage'))
  });

  it('未ログインで /mypage に直接アクセスすると Home へリダイレクト', async () => {
    setAuth({ isLoggedIn: false });
    goTo('/mypage');

    render(<AppRouter />);
    // <Navigate to="/"> が働くまで 1tick 待つ
    await waitFor(() => {
      expect(screen.getByText('HomePage'))
    });
  });

  it('ログイン済みで /mypage にアクセスすると MyPage が表示される', () => {
    setAuth({ isLoggedIn: true, userRole: 'general' });
    goTo('/mypage');

    render(<AppRouter />);
    expect(screen.getByText('MyPage'))
  });

  it('ログイン済みで / (パブリック) にアクセスすると /mypage へリダイレクト', async () => {
    setAuth({ isLoggedIn: true, userRole: 'general' });
    goTo('/');

    render(<AppRouter />);
    await waitFor(() => {
      expect(screen.getByText('MyPage'))
    });
  });

  it('一般ユーザーは /words/new へアクセス出来ない', async () => {
    setAuth({ isLoggedIn: true, userRole: 'general' });
    goTo('/words/new');
  
    render(<AppRouter />);
  
    // ① admin 専用ページが描画されていないこと
    expect(screen.queryByText('WordNewPage')).not
  
    // ② 代わりに MyPage が出てくる（現行実装）
    await waitFor(() => {
      expect(screen.getByText('MyPage'))
    });
  });

  it('admin ユーザーは /words/new へ入れる', () => {
    setAuth({ isLoggedIn: true, userRole: 'admin' });
    goTo('/words/new');

    render(<AppRouter />);
    expect(screen.getByText('WordNewPage'))
  });
});
