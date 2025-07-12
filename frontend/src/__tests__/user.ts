// tests/user.ts  (新規ファイル – 好きな名前でOK)
import userEvent from '@testing-library/user-event';

/**
 * 全テスト共通の userEvent インスタンス
 * ここ**以外で setup() を呼ばない**
 */
export const user = userEvent.setup();
