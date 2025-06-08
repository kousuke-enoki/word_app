/**
 * PublicRoute.test.tsx
 *
 * - vitest / Testing-Library
 * - useAuth をモックして 3 パターン検証
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter, Routes, Route, useLocation } from 'react-router-dom';
import { render, screen } from '@testing-library/react';
import React from 'react';
import PublicRoute from '../PublicRoute';

/* ------------------------------------------------------------------ */
/* ❶ useAuth を好きな戻り値に差し替えられるようモック                  */
/* ------------------------------------------------------------------ */
vi.mock('@/hooks/useAuth', () => ({ useAuth: vi.fn() }));
import { useAuth } from '@/hooks/useAuth';
const mockedUseAuth = vi.mocked(useAuth);

/** 状態切り替え用ヘルパ */
const setAuthState = (state: Partial<ReturnType<typeof useAuth>>) =>
  mockedUseAuth.mockReturnValue({
    isLoggedIn: false,
    isLoading: false,
    userRole: 'guest',
    ...state,
  });

/* ------------------------------------------------------------------ */
/* ❷ 現在パスを可視化するダミーコンポーネント                           */
/* ------------------------------------------------------------------ */
const WhereAmI = () => {
  const loc = useLocation();
  return <p data-testid="path">{loc.pathname}</p>;
};

/* ------------------------------------------------------------------ */
/* ❸ <MemoryRouter> にルートを組んでレンダー                            */
/* ------------------------------------------------------------------ */
const renderWithRouter = (ui: React.ReactElement, start = '/public') =>
  render(
    <MemoryRouter initialEntries={[start]}>
      <Routes>
        <Route path="/mypage" element={<WhereAmI />} />
        <Route path="/public" element={ui} />
      </Routes>
    </MemoryRouter>,
  );

/* ------------------------------------------------------------------ */
/*                               TESTS                                */
/* ------------------------------------------------------------------ */
describe('PublicRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('ロード中は Loading... を表示', () => {
    setAuthState({ isLoading: true });

    renderWithRouter(
      <PublicRoute>
        <p>guest page</p>
      </PublicRoute>,
    );

    expect(screen.getByText('Loading...'))
  });

  it('ログイン済みなら /mypage へリダイレクト', () => {
    setAuthState({ isLoggedIn: true });

    renderWithRouter(
      <PublicRoute>
        <p>guest page</p>
      </PublicRoute>,
    );

    // <Navigate to="/mypage"> が実行された結果、現在パスが /mypage になる
    expect(screen.getByTestId('path').textContent).toBe('/mypage');
  });

  it('未ログインなら子要素を表示', () => {
    setAuthState({ isLoggedIn: false });

    renderWithRouter(
      <PublicRoute>
        <p>welcome guest</p>
      </PublicRoute>,
    );

    expect(screen.getByText('welcome guest'))
  });
});
