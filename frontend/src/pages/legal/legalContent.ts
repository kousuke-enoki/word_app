// src/pages/legal/legalContent.ts
export type Node =
  | { type: 'h1'; text: string }
  | { type: 'h2'; text: string }
  | { type: 'p'; text: string }
  | { type: 'ul'; items: string[] }
  | { type: 'link'; href: string; text: string }
  | { type: 'blockquote'; text: string }
  | { type: 'hr' }

export type LegalDoc = {
  title: string
  updated: string // "YYYY-MM-DD"
  nodes: Node[]
}

const CONTACT_URL = '{{WantedlyのプロフィールURL}}'
const UPDATED = '2025-10-08'

export const legalContent: Record<
  'terms' | 'privacy' | 'cookies' | 'credits',
  LegalDoc
> = {
  terms: {
    title: '利用規約',
    updated: UPDATED,
    nodes: [
      { type: 'h2', text: '1. 本サービスの位置づけ' },
      {
        type: 'p',
        text: '本サービスは、開発者が転職活動のポートフォリオとして作成したデモ版です。商用提供ではなく、予告なく内容の変更・中断・終了を行う場合があります。',
      },
      { type: 'h2', text: '2. サービス内容（デモ版）' },
      {
        type: 'ul',
        items: [
          '英単語データの参照・検索・学習テスト（サンプルデータ中心）',
          '一部機能（登録・編集など）は公開デモでは制限される場合があります',
          'LINEログイン機能は公開時はOFF（機能としては存在）',
        ],
      },
      { type: 'h2', text: '3. アカウント・データ' },
      {
        type: 'ul',
        items: [
          '公開デモはゲスト利用またはサンプルデータ閲覧を想定し、個人を特定できる情報の登録を求めません。',
          '誤って個人情報を送信した場合は、削除依頼をご連絡ください（「連絡先」を参照）。',
        ],
      },
      { type: 'h2', text: '4. 禁止事項' },
      {
        type: 'ul',
        items: [
          '法令・公序良俗に反する行為、権利侵害',
          '過度な負荷・不正アクセス・解析・改変・リバースエンジニアリング',
          '取得データの無断再配布等',
        ],
      },
      { type: 'h2', text: '5. 知的財産' },
      {
        type: 'p',
        text: '本サービスのロゴ・名称・UI等は開発者またはライセンサーに帰属します。オープンソースや辞書データは各ライセンスに従います（/credits参照）。',
      },
      { type: 'h2', text: '6. 免責' },
      {
        type: 'p',
        text: '本サービスは現状有姿で提供します。法令上の責任が認められる場合を除き、本サービスの利用により生じたいかなる損害についても、開発者は責任を負いません。',
      },
      { type: 'h2', text: '7. 規約の変更' },
      {
        type: 'p',
        text: '本ページの更新により効力を生じます。重要な変更は本サイト上で周知します。',
      },
      { type: 'h2', text: '8. 連絡先' },
      {
        type: 'p',
        text: '問い合わせは下記プロフィールのメッセージ機能をご利用ください。',
      },
      { type: 'link', href: CONTACT_URL, text: 'Wantedly プロフィール' },
    ],
  },

  privacy: {
    title: 'プライバシーポリシー',
    updated: UPDATED,
    nodes: [
      { type: 'h2', text: '1. 事業者' },
      { type: 'p', text: '個人開発（ポートフォリオ目的）' },
      { type: 'h2', text: '2. 取得する情報（デモ前提）' },
      {
        type: 'ul',
        items: [
          '公開デモでは、個人を特定可能な情報の収集を行わない構成を基本とします。',
          '操作ログや学習結果は、匿名化されたテストデータまたはサンプルデータとして取り扱います。',
          '誤送信などで個人情報が含まれる場合は削除対応します（「7. 開示・削除」参照）。',
        ],
      },
      { type: 'h2', text: '3. 利用目的' },
      {
        type: 'ul',
        items: [
          '学習機能の動作検証・品質改善（統計化／匿名化）',
          '障害時の原因調査（個人特定を行わない範囲）',
        ],
      },
      { type: 'h2', text: '4. 第三者提供・委託・国外移転' },
      {
        type: 'p',
        text: '公開デモでは第三者提供を行いません。サードパーティ計測・広告・CDN等の外部送信は不使用です（利用開始時に本ページを改定）。',
      },
      { type: 'h2', text: '5. 保存期間' },
      {
        type: 'p',
        text: 'デモ検証用データは最小限・短期間の保存とし、不要になり次第削除します。',
      },
      { type: 'h2', text: '6. 安全管理措置' },
      {
        type: 'ul',
        items: [
          '最小権限、TLS、ログ監査（デモ運用の範囲内）',
          'サンプルデータ中心の構成によるデータ最小化',
        ],
      },
      { type: 'h2', text: '7. 開示・訂正・利用停止・削除' },
      {
        type: 'p',
        text: '誤って送信された個人情報については、合理的範囲で速やかに削除します。連絡先：',
      },
      { type: 'link', href: CONTACT_URL, text: 'Wantedly プロフィール' },
      { type: 'h2', text: '8. 改定' },
      {
        type: 'p',
        text: '本ページの更新により効力を生じます。重要な変更は本サイト上で周知します。',
      },
    ],
  },

  cookies: {
    title: 'クッキーポリシー',
    updated: UPDATED,
    nodes: [
      { type: 'h2', text: '1. 取扱い' },
      {
        type: 'p',
        text: '本サービスの公開デモは、ログイン保持やパフォーマンス計測を行わない構成のため、サードパーティCookieおよび外部送信は使用しません。',
      },
      {
        type: 'p',
        text: '必要に応じてブラウザの LocalStorage 等を機能維持（表示設定など）のみに限定して用いる場合があります。',
      },
      { type: 'h2', text: '2. 管理方法' },
      {
        type: 'ul',
        items: [
          'ブラウザの設定からCookieやサイトデータを削除・無効化できます。',
          '将来、計測や広告を導入する場合は、同意取得とともに本ページを改定し、一覧（ベンダ・送信情報・目的・ポリシー）を掲示します。',
        ],
      },
    ],
  },

  credits: {
    title: 'クレジット／ライセンス表記',
    updated: UPDATED,
    nodes: [
      { type: 'h2', text: '1. 辞書・データソース' },
      {
        type: 'ul',
        items: [
          'JMdict / EDICT — © Electronic Dictionary Research and Development Group (EDRDG) ／ ライセンス：CC BY-SA 3.0（クレジット要件に従い参照利用。再配布は行いません）',
          'WordNet（使用している場合）— © Princeton University ／ ライセンス：Princeton WordNet License',
          'KANJIDIC2（使用している場合）— ライセンス：CC BY-SA 4.0',
        ],
      },
      {
        type: 'blockquote',
        text: '本サービスはアプリ内参照（SaaS的利用）を基本とし、原データのエクスポート等の再配布は行いません。将来、派生データの配布機能を追加する場合は、各ライセンスの継承条件に従い、本ページを改定します。',
      },
      { type: 'h2', text: '2. オープンソース' },
      {
        type: 'ul',
        items: [
          'Frontend：React / Vite / TypeScript ほか',
          'Backend：Go / Gin / ent / PostgreSQL ほか',
          '各プロジェクトのライセンスに従います。',
        ],
      },
      { type: 'h2', text: '3. 商標' },
      {
        type: 'p',
        text: '記載の会社名・製品名・サービス名は、各社の商標または登録商標です。',
      },
    ],
  },
}
